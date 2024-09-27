package utils

import (
	"fmt"
	"strings"

	"osproxy/api/v1alpha3"
	"osproxy/internal/objectStorage"
)

const (
	DefaultSourceKey = "osproxy-source-default"
)

type RequestT struct {
	Host        string
	Port        string
	Path        string
	QueryParams map[string]string
}

func NewRequest(host, fullpath string) (r RequestT) {
	pathQueryParts := strings.SplitN(fullpath, "?", 2)

	r.Path = pathQueryParts[0]
	r.QueryParams = map[string]string{}

	if len(pathQueryParts) == 2 {
		queryParams := strings.Split(pathQueryParts[1], "&")
		for _, qp := range queryParams {
			qpParts := strings.SplitN(qp, "=", 2)
			if len(qpParts) == 2 {
				r.QueryParams[qpParts[0]] = qpParts[1]
			}
		}
	}

	hostParts := strings.Split(host, ":")

	r.Host = hostParts[0]
	if len(hostParts) == 2 {
		r.Port = hostParts[1]
	}

	return r
}

func (req *RequestT) GetObjectFromSource(source v1alpha3.SourceConfigT) (object objectStorage.ObjectT, err error) {
	// Get object path
	originalObjectPath := strings.TrimPrefix(req.Path, "/")

	if source.Type == "host" {
		hostBucketRelation, ok := source.Buckets[req.Host]
		if !ok {
			err = fmt.Errorf("host relation config not provided for '%s' host", req.Host)
			return object, err
		}

		object = setObjectByBucketObject(originalObjectPath, hostBucketRelation)
	}

	if source.Type == "pathPrefix" {
		for prefix, bucketObject := range source.Buckets {
			if strings.HasPrefix(originalObjectPath, prefix) {
				object = setObjectByBucketObject(originalObjectPath, bucketObject)
				break
			}
		}

		if object.Bucket == "" {
			object = setObjectByBucketObject(originalObjectPath, source.Buckets[DefaultSourceKey])
		}
	}

	return object, err
}

func setObjectByBucketObject(objectPath string, bucketObject v1alpha3.BucketObjectConfigT) (object objectStorage.ObjectT) {
	objectPath = strings.TrimPrefix(objectPath, bucketObject.ObjectMod.RemovePrefix)

	if bucketObject.ObjectMod.AddPrefix != "" {
		objectPath = strings.Join([]string{bucketObject.ObjectMod.AddPrefix, objectPath}, "/")
	}

	object = objectStorage.ObjectT{
		Bucket: bucketObject.Bucket,
		Path:   objectPath,
	}

	return object
}

func (r *RequestT) String() string {
	qp := "{"
	for k, v := range r.QueryParams {
		qp += fmt.Sprintf("[%s:%s]", k, v)
	}
	qp += "}"
	return fmt.Sprintf("{host: '%s:%s', path: '%s', queryParams: '%s'}", r.Host, r.Port, r.Path, qp)
}
