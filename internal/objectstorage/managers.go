package objectstorage

import (
	"context"
	"io"
	"osproxy/api/v1alpha5"
)

type ObjectReaderI interface {
	io.Reader
	Close() (err error)
}

type ObjectManagerI interface {
	Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error)
	GetObject(obj ObjectT) (reader ObjectReaderI, info ObjectInfoT, err error)
}

func GetManagers() map[string]ObjectManagerI {
	managers := map[string]ObjectManagerI{
		"s3":  &S3ManagerT{},
		"gcs": &GCSManagerT{},
	}
	return managers
}
