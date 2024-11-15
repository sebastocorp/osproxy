package managers

import (
	"context"
	"fmt"
	"net/http"
	"osproxy/api/v1alpha5"
	"osproxy/internal/objectstorage"
)

type ObjectManagerI interface {
	Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error)
	GetObject(obj objectstorage.ObjectT) (resp *http.Response, err error)
}

func GetManagers() map[string]ObjectManagerI {
	managers := map[string]ObjectManagerI{
		"s3":  &S3ManagerT{},
		"gcs": &GCSManagerT{},
	}
	return managers
}

func GetManager(ctx context.Context, config v1alpha5.ProxySourceConfigT) (manager ObjectManagerI, err error) {
	switch config.Type {
	case "s3":
		{
			manager = &S3ManagerT{}
		}
	case "gcs":
		{
			manager = &GCSManagerT{}
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
