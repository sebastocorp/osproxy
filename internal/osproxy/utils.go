package osproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"osproxy/api/v1alpha2"
	"osproxy/internal/objectStorage"
	"osproxy/internal/utils"

	"gopkg.in/yaml.v3"
)

const (
	defaultRelationKey = "osproxy-default-relation"
)

func (osp *OSProxyT) parseConfig(filepath string) (err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	configBytes = []byte(os.ExpandEnv(string(configBytes)))

	err = yaml.Unmarshal(configBytes, &osp.Config)
	if err != nil {
		return err
	}

	if _, ok := osp.Config.ObjectStorage.Relation.Buckets[defaultRelationKey]; !ok {
		osp.Config.ObjectStorage.Relation.Buckets[defaultRelationKey] = v1alpha2.FrontBackBucketsT{
			Frontend: v1alpha2.BucketSubpathT{
				BucketName: "bucket-frontend-placeholder",
			},
			Backend: v1alpha2.BucketSubpathT{
				BucketName: "bucket-backend-placeholder",
			},
		}
	}

	return err
}

func (osp *OSProxyT) processRequest(r *http.Request) (fObject, bObject objectStorage.ObjectT, req utils.RequestT, err error) {
	// check path
	req = utils.NewRequest(r.Host, r.URL.Path)
	// Get object path
	originalObjectPath := strings.TrimPrefix(req.Path, "/")

	if osp.Config.ObjectStorage.Relation.Type == "host" {
		hostBucketRelation, ok := osp.Config.ObjectStorage.Relation.Buckets[req.Host]
		if !ok {
			err = fmt.Errorf("host relation config not provided for '%s' host", req.Host)
			return fObject, bObject, req, err
		}

		fObject, bObject = osp.setFrontBackBuckets(originalObjectPath, hostBucketRelation)
	}

	if osp.Config.ObjectStorage.Relation.Type == "pathPrefix" {
		for prefix, fbBuckets := range osp.Config.ObjectStorage.Relation.Buckets {
			if strings.HasPrefix(originalObjectPath, prefix) {
				fObject, bObject = osp.setFrontBackBuckets(originalObjectPath, fbBuckets)
				break
			}
		}

		if fObject.BucketName == "" || bObject.BucketName == "" {
			fObject, bObject = osp.setFrontBackBuckets(originalObjectPath, osp.Config.ObjectStorage.Relation.Buckets[defaultRelationKey])
		}
	}

	return fObject, bObject, req, err
}

func (osp *OSProxyT) setFrontBackBuckets(objectPath string, fbBuckets v1alpha2.FrontBackBucketsT) (fObject, bObject objectStorage.ObjectT) {
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

	transfer := TransferT{
		From: bObject,
		To:   fObject,
	}

	bodyBytes, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	http.DefaultClient.Timeout = 100 * time.Millisecond
	resp, err := http.Post(osp.Config.Action.APICall.URL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	resp.Body.Close()

	return err
}
