package proxyComp

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"

	"osproxy/api/v1alpha3"
	"osproxy/internal/logger"
	"osproxy/internal/objectStorage"
	"osproxy/internal/pools"
	"osproxy/internal/utils"
)

const (
	logExtraFieldKeyRequest       = "request"
	logExtraFieldKeyStatusCode    = "status_code"
	logExtraFieldKeyContentLength = "content_length"
	logExtraFieldKeyDataBytes     = "data_bytes"
	logExtraFieldKeyObject        = "object"
	logExtraFieldKeyError         = "error"
)

type ProxyT struct {
	config v1alpha3.ProxyConfigT
	log    logger.LoggerT

	actionPool *pools.ActionPoolT
	objManager objectStorage.ManagerT
}

func NewProxy(config v1alpha3.ProxyConfigT, actionPool *pools.ActionPoolT) (p ProxyT, err error) {
	p.config = config
	p.actionPool = actionPool

	if p.config.Loglevel == "" {
		p.config.Loglevel = "info"
	}

	level, err := logger.GetLevel(p.config.Loglevel)
	if err != nil {
		log.Fatalf("unable to get log level in proxy config: %s", err.Error())
	}

	p.log = logger.NewLogger(context.Background(), level, map[string]any{
		"service":   "osproxy",
		"component": "proxy",
	})

	p.objManager, err = objectStorage.NewManager(context.Background(),
		p.config.Source.Config,
	)

	return p, err
}
func (p *ProxyT) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	err := http.ListenAndServe(
		fmt.Sprintf("%s:%s", p.config.Address, p.config.Port),
		http.HandlerFunc(p.HandleFunc),
	)
	if err != nil {
		p.log.Error("unable to serve proxy",
			map[string]any{
				"error": err.Error(),
			},
		)
	}
}

func (p *ProxyT) HandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	logExtraFields := map[string]any{
		logExtraFieldKeyRequest:       "none",
		logExtraFieldKeyStatusCode:    "none",
		logExtraFieldKeyDataBytes:     "none",
		logExtraFieldKeyContentLength: "none",
		logExtraFieldKeyObject:        "none",
		logExtraFieldKeyError:         "none",
	}

	req := utils.NewRequest(r.Host, r.URL.Path)
	logExtraFields[logExtraFieldKeyRequest] = req.String()

	object, err := req.GetObjectFromSource(p.config.Source)
	if err != nil {
		logExtraFields[logExtraFieldKeyError] = err.Error()
		logExtraFields[logExtraFieldKeyStatusCode] = http.StatusInternalServerError
		p.log.Debug("unable to process request", logExtraFields)
		p.requestResponseErrorLog(w, http.StatusInternalServerError, "Internal Server Error", "unable to handle request", logExtraFields)
		return
	}
	p.log.Debug("success in process request", logExtraFields)

	logExtraFields[logExtraFieldKeyObject] = object.String()

	objectResp, info, err := p.objManager.S3GetObject(object)
	if err != nil {
		logExtraFields[logExtraFieldKeyError] = err.Error()
		logExtraFields[logExtraFieldKeyStatusCode] = http.StatusInternalServerError
		p.log.Debug("unable to get object", logExtraFields)
		p.requestResponseErrorLog(w, http.StatusInternalServerError, "Internal Server Error", "unable to handle request", logExtraFields)
		return
	}
	defer objectResp.Close()
	p.log.Debug("success in get object", logExtraFields)

	if !info.Exist {
		logExtraFields[logExtraFieldKeyError] = "object not exist in bucket"
		logExtraFields[logExtraFieldKeyStatusCode] = http.StatusInternalServerError
		p.requestResponseErrorLog(w, http.StatusNotFound, "Not Found", "unable to handle request", logExtraFields)

		p.log.Debug("add action in pool", logExtraFields)
		p.actionPool.Add(pools.ActionPoolRequestT{
			Object: object,
		})
		return
	}

	// Set headers before response body
	contentLen := strconv.FormatInt(info.Size, 10)
	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("Content-Length", contentLen)
	if filename, ok := req.QueryParams["filename"]; ok {
		contentDispositionHeaderVal := fmt.Sprintf("inline; filename=\"%s\"", filename)
		w.Header().Set("Content-Disposition", contentDispositionHeaderVal)
	}

	w.WriteHeader(http.StatusOK)

	logExtraFields[logExtraFieldKeyStatusCode] = http.StatusOK
	logExtraFields[logExtraFieldKeyContentLength] = contentLen
	// Copy object data in response body
	dataBytes, dataErr := io.Copy(w, objectResp)
	logExtraFields[logExtraFieldKeyDataBytes] = dataBytes
	if dataErr != nil {
		logExtraFields[logExtraFieldKeyError] = dataErr.Error()
		p.log.Debug("unable to copy data", logExtraFields)
		p.requestResponseErrorLog(w, http.StatusInternalServerError, "Internal Server Error", "unable to handle request", logExtraFields)
		return
	}

	p.log.Info("success in handle request", logExtraFields)
}
