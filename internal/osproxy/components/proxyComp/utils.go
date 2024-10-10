package proxyComp

import (
	"fmt"
	"net/http"
	"osproxy/api/v1alpha4"
	"osproxy/internal/objectStorage"
	"strconv"
	"strings"
)

func (p *ProxyT) requestResponseError(respWriter http.ResponseWriter, respStatusCode int, respMessage string) {
	respMessage = fmt.Sprintf("%d %s\n", respStatusCode, respMessage)

	// response to user request
	respWriter.Header().Set("Content-Type", "text/plain")
	respWriter.Header().Set("Content-Length", strconv.Itoa(len(respMessage)))
	respWriter.WriteHeader(respStatusCode)
	respWriter.Write([]byte(respMessage))
}

func (p *ProxyT) GetObjectFromRequest(r *http.Request) (object objectStorage.ObjectT, err error) {
	// Get object path
	originalObjectPath := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "?")[0]

	var mod v1alpha4.ObjectModificationConfigT
	var found bool = false
	switch p.config.RequestRouting.Type {
	case "host":
		{
			mod, found = p.config.RequestRouting.Routes[r.Host]
		}
	case "headerValue":
		{
			if p.config.RequestRouting.HeaderName == "" {
				err = fmt.Errorf("header name in header routing config not provided")
				return object, err
			}

			mod, found = p.config.RequestRouting.Routes[r.Header.Get(p.config.RequestRouting.HeaderName)]

		}
	case "pathPrefix":
		{
			for prefix, objMod := range p.config.RequestRouting.Routes {
				if strings.HasPrefix(originalObjectPath, prefix) {
					mod = objMod
					found = true
					break
				}
			}
		}
	}

	if !found {
		err = fmt.Errorf("routing config not provided for this request")
		return object, err
	}

	object = setObjectByBucketObject(originalObjectPath, mod)

	return object, err
}

func setObjectByBucketObject(objectPath string, objectMod v1alpha4.ObjectModificationConfigT) (object objectStorage.ObjectT) {
	objectPath = strings.TrimPrefix(objectPath, objectMod.RemovePrefix)
	objectPath = objectMod.AddPrefix + objectPath

	object = objectStorage.ObjectT{
		Bucket: objectMod.Bucket,
		Path:   objectPath,
	}

	return object
}
