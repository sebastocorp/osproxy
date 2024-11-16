package proxyComp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"osproxy/api/v1alpha5"
	"osproxy/internal/global"
	"osproxy/internal/logger"
	"osproxy/internal/pools"
	"osproxy/internal/sources/managers"
	"osproxy/internal/utils"
)

type ProxyT struct {
	ctx    context.Context
	log    logger.LoggerT
	config *v1alpha5.OSProxyConfigT

	server     *http.Server
	actionPool *pools.ActionPoolT

	sources          map[string]managers.ObjectManagerI
	requestModifiers map[string]*v1alpha5.ProxyModifierConfigT
}

func NewProxy(config *v1alpha5.OSProxyConfigT, actionPool *pools.ActionPoolT) (p ProxyT, err error) {
	p.config = config
	p.actionPool = actionPool
	p.ctx = context.Background()

	logCommon := global.GetLogCommonFields()
	logCommon[global.LogFieldKeyCommonComponent] = global.LogFieldValueComponentProxy
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

	err := p.server.ListenAndServe()
	if err != nil {
		logExtra[global.LogFieldKeyExtraError] = err.Error()
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
	logExtraFields[global.LogFieldKeyExtraRequest] = utils.RequestString(req)

	route, err := p.getRouteFromRequest(req)
	if err != nil {
		p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")

		logExtraFields[global.LogFieldKeyExtraError] = err.Error()
		logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusInternalServerError
		p.log.Error("unable to process request", logExtraFields)
		return
	}

	p.modRequest(req, route.Modifiers)

	resp, err := p.sources[route.Source].GetObject(req, route.Bucket)
	if err != nil {
		logExtraFields[global.LogFieldKeyExtraError] = err.Error()
		logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusInternalServerError
		p.log.Error("unable to get object", logExtraFields)
		p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	p.log.Debug("success in get object reader", logExtraFields)

	for reaci, reacv := range p.config.Proxy.RespReactions {
		_ = reaci
		switch reacv.Condition.Key {
		case ":host:":
			{
				if resp.Request.Host == reacv.Condition.Value {
					p.log.Error("respopnse reaction with host not implemented", logExtraFields)
				}
			}
		case ":status:":
			{
				if reacv.Condition.Value == strconv.Itoa(resp.StatusCode) {
					if reacv.Type == "ResponseSustitution" {
						resp.Body.Close()
						resp, err = p.sources[reacv.ResponseSustitution.Source].GetObject(r, route.Bucket)
						if err != nil {
							logExtraFields[global.LogFieldKeyExtraError] = err.Error()
							logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusInternalServerError
							p.log.Error("unable to get object", logExtraFields)
							p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")
							return
						}
						defer resp.Body.Close()
					}
				}
			}
		default:
			{
				headerValue := resp.Header.Get(reacv.Condition.Key)
				if headerValue == reacv.Condition.Value {
					p.log.Error("respopnse reaction with headers not implemented", logExtraFields)
				}
			}
		}
	}

	// Set headers before response body
	for hk, hvs := range resp.Header {
		for _, hv := range hvs {
			w.Header().Set(hk, hv)
		}
	}
	w.Header().Set("Connection", "close")

	w.WriteHeader(http.StatusOK)

	logExtraFields[global.LogFieldKeyExtraResponse] = utils.ResponseString(resp)
	// Copy object data in response body
	dataBytes, dataErr := io.Copy(w, resp.Body)
	logExtraFields[global.LogFieldKeyExtraDataBytes] = dataBytes
	if dataErr != nil {
		logExtraFields[global.LogFieldKeyExtraError] = dataErr.Error()
		p.log.Error("unable to copy data", logExtraFields)
		p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	p.log.Info("success in handle request", logExtraFields)
}
