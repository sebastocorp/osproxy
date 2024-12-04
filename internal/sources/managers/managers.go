package managers

import (
	"context"
	"fmt"
	"net/http"

	"osproxy/api/v1alpha5"
)

type ObjectManagerI interface {
	Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error)
	GetObject(r *http.Request, bucket string) (resp *http.Response, err error)
}

func GetManager(ctx context.Context, config v1alpha5.ProxySourceConfigT) (manager ObjectManagerI, err error) {
	switch config.Type {
	case "S3":
		{
			manager = &S3ManagerT{}
		}
	case "GCS":
		{
			manager = &GCSManagerT{}
		}
	case "HTTP":
		{
			manager = &HTTPManagerT{}
		}
	default:
		{
			err = fmt.Errorf("unsuported manager type")
			return manager, err
		}
	}
	err = manager.Init(ctx, config)
	return manager, err
}
