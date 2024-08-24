package osproxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"osproxy/api/v1alpha1"
	"osproxy/internal/logger"
	"osproxy/internal/objectStorage"

	"github.com/minio/minio-go/v7"
)

type OSProxyT struct {
	config     v1alpha1.OSProxyConfigT
	objManager objectStorage.ManagerT
}

func NewOSProxy(config string) (osp OSProxyT, err error) {
	err = osp.parseConfig(config)
	if err != nil {
		return osp, err
	}

	osp.objManager, err = objectStorage.NewManager(context.Background(),
		objectStorage.S3T{
			Endpoint:        osp.config.OSConfig.S3.Endpoint,
			AccessKeyID:     osp.config.OSConfig.S3.AccessKeyID,
			SecretAccessKey: osp.config.OSConfig.S3.SecretAccessKey,
		},
		objectStorage.GCST{
			CredentialsFile: osp.config.OSConfig.GCS.CredentialsFile,
		},
	)

	return osp, err
}

func (osp *OSProxyT) HandleFunc(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
		}
	}()

	logger.Log.Infof("handle request '%s'", r.URL.Path)
	fObject, bObject, req, err := osp.processRequest(r)
	if err != nil {
		logger.Log.Errorf("unable to process request '%s': %s", r.URL.Path, err.Error())
		return
	}
	_ = bObject

	logger.Log.Infof("check object '%s'", fObject.String())
	exist, info, err := osp.objManager.S3ObjectExist(fObject)
	if err != nil {
		logger.Log.Errorf("unable to check object %s: %s", fObject.String(), err.Error())
		return
	}

	if !exist {
		logger.Log.Errorf("object %s not exist, making transfer request", fObject.String())
		err = osp.makeAPICall(fObject, bObject)
		if err != nil {
			logger.Log.Errorf("unable to request transfer %s to %s: %s",
				bObject.String(), fObject.String(), err.Error())
		}
		err = fmt.Errorf("object NOT exist")
		return
	}

	logger.Log.Infof("get object %s", fObject.String())
	object, err := osp.objManager.S3.Client.GetObject(osp.objManager.Ctx, fObject.BucketName, fObject.ObjectPath, minio.GetObjectOptions{})
	if err != nil {
		logger.Log.Errorf("unable get object '%s': %s", fObject.String(), err.Error())
		return
	}
	defer object.Close()

	// Set headers before response body
	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size, 10))
	if filename, ok := req.QueryParams["filename"]; ok {
		contentDispositionHeaderVal := fmt.Sprintf("inline; filename=\"%s\"", filename)
		w.Header().Set("Content-Disposition", contentDispositionHeaderVal)
	}

	// Copy object data in response body
	if _, err := io.Copy(w, object); err != nil {
		logger.Log.Errorf("unable to set object %s in response body: %s", fObject.String(), err.Error())
	}
}
