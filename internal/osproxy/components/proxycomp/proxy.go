package proxycomp

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"osproxy/api/v1alpha5"
	"osproxy/internal/global"
	"osproxy/internal/logger"
	"osproxy/internal/sources"
	"osproxy/internal/sources/managers"
	"osproxy/internal/utils"
)

type ProxyT struct {
	ctx    context.Context
	log    logger.LoggerT
	config *v1alpha5.OSProxyConfigT

	server *http.Server

	sources          map[string]managers.ObjectManagerI
	requestModifiers map[string]*v1alpha5.ProxyModifierConfigT
}

func NewProxy(config *v1alpha5.OSProxyConfigT) (p ProxyT, err error) {
	p.config = config
	p.ctx = context.Background()

	logCommon := global.GetLogCommonFields()
	global.SetLogExtraField(logCommon, global.LogFieldKeyCommonComponent, global.LogFieldValueComponentProxy)
	p.log = logger.NewLogger(p.ctx, logger.GetLevel(p.config.Proxy.Loglevel), logCommon)

	mux := http.NewServeMux()
	mux.HandleFunc(global.EndpointHealthz, p.getHealthz)
	mux.HandleFunc("/", p.HandleFunc)

	p.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", p.config.Proxy.Address, p.config.Proxy.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	p.requestModifiers = map[string]*v1alpha5.ProxyModifierConfigT{}
	for modi, modv := range p.config.Proxy.RequestModifiers {
		p.requestModifiers[modv.Name] = &p.config.Proxy.RequestModifiers[modi]
	}

	p.sources = map[string]managers.ObjectManagerI{}
	for _, srcv := range p.config.Proxy.Sources {
		p.sources[srcv.Name], err = managers.GetManager(p.ctx, srcv)
		if err != nil {
			return p, err
		}
	}

	return p, err
}
func (p *ProxyT) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	logExtra := global.GetLogExtraFieldsProxy()

	p.log.Info("proxy initialized", logExtra)
	err := p.server.ListenAndServe()
	if err != nil {
		global.SetLogExtraField(logExtra, global.LogFieldKeyExtraError, err.Error())
		p.log.Error("unable to serve proxy", logExtra)
	}
}

func (p *ProxyT) getHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (p *ProxyT) HandleFunc(w http.ResponseWriter, r *http.Request) {
	req := r.Clone(p.ctx)
	defer r.Body.Close()
	defer req.Body.Close()

	var err error
	logExtraFields := global.GetLogExtraFieldsProxy()
	reqStr := utils.RequestString(req)
	md5Hash := md5.New()
	_, err = md5Hash.Write([]byte(reqStr))
	if err != nil {
		logResp := p.requestResponseError(w, http.StatusInternalServerError)

		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraError, err.Error())
		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraResponse, logResp)
		p.log.Error("unable to request id hash", logExtraFields)
		return
	}
	global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraRequest, reqStr)
	global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraRequestId, hex.EncodeToString(md5Hash.Sum(nil)))
	p.log.Info("handle request", logExtraFields)
	global.ResetLogExtraField(logExtraFields, global.LogFieldKeyExtraRequest)

	route, err := p.getRouteFromRequest(req)
	if err != nil {
		logResp := p.requestResponseError(w, http.StatusInternalServerError)

		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraError, err.Error())
		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraResponse, logResp)
		p.log.Error("unable to get route from request", logExtraFields)
		return
	}

	p.modRequest(req, route.Modifiers)

	resp, err := p.sources[route.Source].GetObject(req, route.Bucket)
	if err != nil {
		logResp := p.requestResponseError(w, http.StatusInternalServerError)

		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraError, err.Error())
		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraResponse, logResp)
		p.log.Error("unable to make source request", logExtraFields)
		return
	}
	defer resp.Body.Close()

	p.log.Debug("success in source request", logExtraFields)

	for _, reacv := range p.config.Proxy.RespReactions {
		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraReaction, fmt.Sprintf("{name: '%s', type: '%s'}", reacv.Name, reacv.Type))
		p.log.Info("execute reaction", logExtraFields)
		switch reacv.Condition.Key {
		case ":host:":
			{
				if resp.Request.Host == reacv.Condition.Value {
					p.log.Error("response reaction with host not implemented", logExtraFields)
				}
			}
		case ":status:":
			{
				if reacv.Condition.Value == strconv.Itoa(resp.StatusCode) {
					switch reacv.Type {
					case "ResponseSustitution":
						{
							resp2, err := p.sources[reacv.ResponseSustitution.Source].GetObject(r, route.Bucket)
							if err != nil {
								global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraError, err.Error())
								p.log.Error("unable to make source request", logExtraFields)
								global.ResetLogExtraField(logExtraFields, global.LogFieldKeyExtraError)
								continue
							}
							resp.Body.Close()
							resp = resp2
							defer resp2.Body.Close()
						}
					case "PostObject":
						{
							object := sources.ObjectT{
								Bucket: route.Bucket,
								Path:   strings.TrimPrefix(req.URL.Path, "/"),
							}

							for hk := range req.Header {
								object.Metadata.Set(hk, req.Header.Get(hk))
							}

							data, err := json.Marshal(object)
							if err != nil {
								global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraError, err.Error())
								p.log.Error("unable to post object json", logExtraFields)
								global.ResetLogExtraField(logExtraFields, global.LogFieldKeyExtraError)
								continue
							}

							http.DefaultClient.Timeout = 5 * time.Second
							respPost, err := http.Post(reacv.PostObject.Endpoint, "application/json", bytes.NewBuffer(data))
							if err != nil {
								global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraError, err.Error())
								p.log.Error("unable to post object json", logExtraFields)
								global.ResetLogExtraField(logExtraFields, global.LogFieldKeyExtraError)
								continue
							}
							respPost.Body.Close()
						}
					}
				}
			}
		default:
			{
				headerValue := resp.Header.Get(reacv.Condition.Key)
				if headerValue == reacv.Condition.Value {
					p.log.Error("response reaction with headers not implemented", logExtraFields)
				}
			}
		}
	}
	global.ResetLogExtraField(logExtraFields, global.LogFieldKeyExtraReaction)

	// Set headers before response body
	for hk, hvs := range resp.Header {
		for _, hv := range hvs {
			w.Header().Set(hk, hv)
		}
	}
	w.Header().Set("Connection", "close")

	w.WriteHeader(resp.StatusCode)

	// Copy object data in response body
	dataBytes, dataErr := io.Copy(w, resp.Body)
	global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraDataBytes, dataBytes)
	if dataErr != nil {
		global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraError, dataErr.Error())
		p.log.Error("unable to copy data", logExtraFields)
		return
	}

	global.SetLogExtraField(logExtraFields, global.LogFieldKeyExtraResponse, utils.ResponseString(resp))
	p.log.Info("success in handle request", logExtraFields)
}
