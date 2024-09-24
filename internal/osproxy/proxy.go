package osproxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"osproxy/api/v1alpha2"
	"osproxy/internal/logger"
	"osproxy/internal/objectStorage"
)

type OSProxyT struct {
	Config     v1alpha2.OSProxyConfigT
	objManager objectStorage.ManagerT
}

func NewOSProxy(config string) (osp OSProxyT, err error) {
	err = osp.parseConfig(config)
	if err != nil {
		return osp, err
	}

	osp.objManager, err = objectStorage.NewManager(context.Background(),
		osp.Config.ObjectStorage.S3,
		osp.Config.ObjectStorage.GCS,
	)

	return osp, err
}

func (osp *OSProxyT) HandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeDirectResponse(w, http.StatusInternalServerError, "Internal Server Error")
		}
	}()

	fObject, bObject, req, err := osp.processRequest(r)
	if err != nil {
		logger.Log.Errorf([]any{"status_code", http.StatusInternalServerError}, "unable to process request '%s': %s", r.URL.Path, err.Error())
		return
	}
	extra := []any{
		"request", req.String(),
		"object", fObject.String(),
		"backend_object", bObject.String(),
	}
	logger.Log.Infof(extra, "handle request")

	logger.Log.Infof(extra, "get object")
	object, info, err := osp.objManager.S3GetObject(fObject)
	if err != nil {
		extra = append(extra, "error", err.Error(), "status_code", http.StatusInternalServerError)
		logger.Log.Errorf(extra, "unable to get object")
		return
	}
	defer object.Close()

	if !info.Exist {
		writeDirectResponse(w, http.StatusNotFound, "Not Found")
		extra = append(extra, "status_code", http.StatusNotFound)

		logger.Log.Errorf(extra, "object not exist, making actions")
		err = osp.makeAPICall(fObject, bObject)
		if err != nil {
			extraActions := append(extra, "error", err.Error())
			logger.Log.Errorf(extraActions, "unable to execute actions")
			err = nil
		} else {
			logger.Log.Infof(extra, "success executing actions")
		}
		return
	}

	extra = append(extra, "status_code", http.StatusOK)

	// Set headers before response body
	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
	if filename, ok := req.QueryParams["filename"]; ok {
		contentDispositionHeaderVal := fmt.Sprintf("inline; filename=\"%s\"", filename)
		w.Header().Set("Content-Disposition", contentDispositionHeaderVal)
	}

	w.WriteHeader(http.StatusOK)

	// Copy object data in response body
	if _, dataErr := io.Copy(w, object); dataErr != nil {
		extra = append(extra, "error", err.Error())
		logger.Log.Errorf(extra, "unable to set data in response body")
	} else {
		logger.Log.Infof(extra, "success handling request")
	}
}
