package proxyComp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"osproxy/api/v1alpha3"
	"osproxy/internal/global"
	"osproxy/internal/logger"
	"osproxy/internal/objectStorage"
	"osproxy/internal/pools"
	"osproxy/internal/utils"
)

type ProxyT struct {
	config v1alpha3.ProxyConfigT
	log    logger.LoggerT

	actionPool *pools.ActionPoolT
	ctx        context.Context
	server     *http.Server
	objManager objectStorage.ManagerT
}

func NewProxy(config v1alpha3.ProxyConfigT, actionPool *pools.ActionPoolT) (p ProxyT, err error) {
	p.config = config
	p.actionPool = actionPool

	logCommon := global.GetLogCommonFields()
	logCommon[global.LogFieldKeyCommonComponent] = global.LogFieldValueComponentProxy
	p.log = logger.NewLogger(context.Background(), logger.GetLevel(p.config.Loglevel), logCommon)

	mux := http.NewServeMux()
	mux.HandleFunc(global.EndpointHealthz, p.getHealthz)
	mux.HandleFunc("/", p.HandleFunc)

	p.ctx = context.Background()
	p.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", p.config.Address, p.config.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	p.objManager, err = objectStorage.NewManager(context.Background(),
		p.config.Source.Config,
	)

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
	var err error
	logExtraFields := global.GetLogExtraFieldsProxy()

	req := utils.NewRequest(r.Host, r.URL.Path)
	logExtraFields[global.LogFieldKeyExtraRequest] = req.String()

	object, err := req.GetObjectFromSource(p.config.Source)
	if err != nil {
		p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")

		logExtraFields[global.LogFieldKeyExtraError] = err.Error()
		logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusInternalServerError
		p.log.Error("unable to process request", logExtraFields)
		return
	}
	p.log.Debug("success in process request", logExtraFields)

	logExtraFields[global.LogFieldKeyExtraObject] = object.String()

	objectResp, info, err := p.objManager.S3GetObject(object)
	if err != nil {
		logExtraFields[global.LogFieldKeyExtraError] = err.Error()
		if info.NotExistError {
			p.requestResponseError(w, http.StatusNotFound, "Not Found")

			logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusNotFound
			p.log.Error("object does not exist", logExtraFields)

			p.actionPool.Add(pools.ActionPoolRequestT{
				Object: object,
			})
			p.log.Debug("add action in pool", logExtraFields)

			return
		}

		logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusInternalServerError
		p.log.Error("unable to get object", logExtraFields)
		p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer objectResp.Close()
	p.log.Debug("success in get object", logExtraFields)

	// Set headers before response body
	contentLen := strconv.FormatInt(info.Size, 10)
	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("Content-Length", contentLen)
	if filename, ok := req.QueryParams["filename"]; ok {
		contentDispositionHeaderVal := fmt.Sprintf("inline; filename=\"%s\"", filename)
		w.Header().Set("Content-Disposition", contentDispositionHeaderVal)
	}

	w.WriteHeader(http.StatusOK)

	logExtraFields[global.LogFieldKeyExtraStatusCode] = http.StatusOK
	logExtraFields[global.LogFieldKeyExtraContentLength] = contentLen
	// Copy object data in response body
	dataBytes, dataErr := io.Copy(w, objectResp)
	logExtraFields[global.LogFieldKeyExtraDataBytes] = dataBytes
	if dataErr != nil {
		logExtraFields[global.LogFieldKeyExtraError] = dataErr.Error()
		p.log.Error("unable to copy data", logExtraFields)
		p.requestResponseError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	p.log.Info("success in handle request", logExtraFields)
}
