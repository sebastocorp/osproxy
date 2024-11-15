package proxyComp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"osproxy/api/v1alpha5"
	"osproxy/internal/global"
	"osproxy/internal/logger"
	"osproxy/internal/modifiers"
	"osproxy/internal/objectstorage"
	"osproxy/internal/objectstorage/managers"
	"osproxy/internal/pools"
	"osproxy/internal/utils"
)

type ProxyT struct {
	ctx    context.Context
	log    logger.LoggerT
	config *v1alpha5.OSProxyConfigT

	server     *http.Server
	actionPool *pools.ActionPoolT

	sources map[string]managers.ObjectManagerI
	routes  map[string]routeT
}

type routeT struct {
	source    *managers.ObjectManagerI
	modifiers []modifiers.ModifierT
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
	defer r.Body.Close()
	req := r.Clone(p.ctx)
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
	object := objectstorage.ObjectT{
		Bucket: route.Bucket,
		Path:   req.URL.Path,
	}
	logExtraFields[global.LogFieldKeyExtraObject] = object.String()

	resp, err := p.sources[route.Source].GetObject(object)
	if err != nil {
		logExtraFields[global.LogFieldKeyExtraError] = err.Error()
		logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusInternalServerError
		p.log.Error("unable to get object", logExtraFields)
		p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	p.log.Debug("success in get object reader", logExtraFields)

	// Set headers before response body
	w.WriteHeader(http.StatusOK)
	for hk, hvs := range resp.Header {
		for _, hv := range hvs {
			w.Header().Set(hk, hv)
		}
	}

	logExtraFields[global.LogFieldKeyExtraStatusCode] = resp.StatusCode
	logExtraFields[global.LogFieldKeyExtraContentLength] = resp.Header.Get("Content-Length")
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
