package osproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"osproxy/api/v1alpha1"
	"osproxy/internal/config"
	"osproxy/internal/objectStorage"
	"osproxy/internal/utils"
	"strings"
)

func (osp *OSProxyT) parseConfig(filepath string) (err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	configBytes = []byte(os.ExpandEnv(string(configBytes)))

	osp.config, err = config.Parse(configBytes)

	return err
}

func (osp *OSProxyT) processRequest(r *http.Request) (fObject, bObject objectStorage.ObjectT, req utils.RequestT, err error) {
	// check path
	req = utils.NewRequest(r.Host, r.URL.Path)
	// Get object path
	originalObjectPath := strings.TrimPrefix(req.Path, "/")

	if osp.config.Relation.Type == "host" {
		hostBucketRelation, ok := osp.config.Relation.Buckets[req.Host]
		if !ok {
			err = fmt.Errorf("host relation config not provided for '%s' host", req.Host)
			return fObject, bObject, req, err
		}

		fObject, bObject = osp.setFrontBackBuckets(originalObjectPath, hostBucketRelation)
	}

	if osp.config.Relation.Type == "pathPrefix" {
		for prefix, fbBuckets := range osp.config.Relation.Buckets {
			if strings.HasPrefix(originalObjectPath, prefix) {
				fObject, bObject = osp.setFrontBackBuckets(originalObjectPath, fbBuckets)
				break
			}
		}
	}

	return fObject, bObject, req, err
}

func (osp *OSProxyT) setFrontBackBuckets(objectPath string, fbBuckets v1alpha1.FrontBackBucketsT) (fObject, bObject objectStorage.ObjectT) {
	fObject = objectStorage.ObjectT{
		BucketName: fbBuckets.Frontend.BucketName,
		ObjectPath: objectPath,
	}
	if fbBuckets.Frontend.AddPathPrefix != "" {
		fObject.ObjectPath = strings.Join([]string{fbBuckets.Frontend.AddPathPrefix, fObject.ObjectPath}, "/")
	}
	if fbBuckets.Frontend.RemovePathPrefix != "" {
		fObject.ObjectPath = strings.TrimPrefix(fObject.ObjectPath, fbBuckets.Frontend.RemovePathPrefix)
	}

	bObject = objectStorage.ObjectT{
		BucketName: fbBuckets.Backend.BucketName,
		ObjectPath: objectPath,
	}
	if fbBuckets.Backend.AddPathPrefix != "" {
		bObject.ObjectPath = strings.Join([]string{fbBuckets.Backend.AddPathPrefix, bObject.ObjectPath}, "/")
	}
	if fbBuckets.Backend.RemovePathPrefix != "" {
		bObject.ObjectPath = strings.TrimPrefix(bObject.ObjectPath, fbBuckets.Backend.RemovePathPrefix)
	}

	return fObject, bObject
}

func (osp *OSProxyT) makeAPICall(fObject, bObject objectStorage.ObjectT) (err error) {
	type TransferT struct {
		From objectStorage.ObjectT `json:"from"`
		To   objectStorage.ObjectT `json:"to"`
	}

	type apiTransferRequestT struct {
		Transfer TransferT `json:"transfer"`
	}

	body := apiTransferRequestT{
		Transfer: TransferT{
			From: bObject,
			To:   fObject,
		},
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s:%s%s",
		osp.config.TransferService.Host,
		osp.config.TransferService.Port,
		osp.config.TransferService.Endpoint,
	)
	_, err = http.Post(requestURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	return err
}
