package routing

import (
	"context"
	"fmt"
	"net/http"
	"osproxy/api/v1alpha5"
	"osproxy/internal/objectstorage"
)

type RequestRouteI interface {
	Init(ctx context.Context, config v1alpha5.ProxyRequestRoutingConfigT) (err error)
	GetObject(obj objectstorage.ObjectT) (resp *http.Response, err error)
}

func GetManager(ctx context.Context, config v1alpha5.ProxyRequestRoutingConfigT) (router RequestRouteI, err error) {
	switch config.MatchType {
	case "host":
		{
		}
	case "headerValue":
		{
		}
	case "pathPrefix":
		{
		}
	default:
		{
			err = fmt.Errorf("unsuported router type")
			return router, err
		}
	}

	return router, err
}
