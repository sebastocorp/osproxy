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
	statusCode := http.StatusInternalServerError
	defer func() {
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		} else if statusCode == http.StatusNotFound {
			http.Error(w, "not found", statusCode)
		}
	}()

	logger.Log.Infof("handle request '%s'", r.URL.Path)
	fObject, bObject, req, err := osp.processRequest(r)
	if err != nil {
		logger.Log.Errorf("unable to process request '%s': %s", r.URL.Path, err.Error())
		return
	}

	logger.Log.Infof("get object '%s'", fObject.String())
	object, info, err := osp.objManager.S3GetObject(fObject)
	if err != nil {
		logger.Log.Errorf("unable to get object %s: %s", fObject.String(), err.Error())
		return
	}
	defer object.Close()

	if !info.Exist {
		statusCode = http.StatusNotFound

		logger.Log.Errorf("object %s not exist, making actions", fObject.String())
		err = osp.makeAPICall(fObject, bObject)
		if err != nil {
			logger.Log.Errorf("unable to make transfer request from %s to %s: %s",
				bObject.String(), fObject.String(), err.Error())
			err = nil
		}
		return
	}

	// Set headers before response body
	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
	if filename, ok := req.QueryParams["filename"]; ok {
		contentDispositionHeaderVal := fmt.Sprintf("inline; filename=\"%s\"", filename)
		w.Header().Set("Content-Disposition", contentDispositionHeaderVal)
	}

	w.WriteHeader(http.StatusOK)

	// Copy object data in response body
	if _, err := io.Copy(w, object); err != nil {
		logger.Log.Errorf("unable to set object %s in response body: %s", fObject.String(), err.Error())
	}
}
