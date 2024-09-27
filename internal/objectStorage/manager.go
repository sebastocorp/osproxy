package objectStorage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"

	"osproxy/api/v1alpha3"

	"cloud.google.com/go/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/api/option"
)

type ManagerT struct {
	Ctx       context.Context
	S3Client  *minio.Client
	GCSClient *storage.Client
}

type ObjectT struct {
	Bucket string      `json:"bucket"`
	Path   string      `json:"path"`
	Info   ObjectInfoT `json:"-"`
}

type ObjectInfoT struct {
	NotExistError bool
	MD5           string
	Size          int64
	ContentType   string
}

func NewManager(ctx context.Context, config v1alpha3.ObjectStorageConfigT) (man ManagerT, err error) {
	man.Ctx = ctx

	man.S3Client, err = minio.New(
		config.S3.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.S3.AccessKeyID, config.S3.SecretAccessKey, ""),
			Region: config.S3.Region,
			Secure: config.S3.Secure,
		},
	)
	if err != nil {
		return man, err
	}

	if reflect.ValueOf(config.S3).IsZero() {
		man.GCSClient, err = storage.NewClient(man.Ctx, option.WithCredentialsFile(config.GCS.CredentialsFile))
	}

	return man, err
}

func (m *ManagerT) S3GetObject(obj ObjectT) (object *minio.Object, info ObjectInfoT, err error) {
	object, err = m.S3Client.GetObject(context.Background(), obj.Bucket, obj.Path, minio.GetObjectOptions{})
	if err != nil {
		return object, info, err
	}

	stat, err := object.Stat()
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			info.NotExistError = true
		}
		return object, info, err
	}

	info.NotExistError = false
	info.MD5 = stat.ETag
	info.ContentType = stat.ContentType
	info.Size = stat.Size

	return object, info, err
}

func (o *ObjectT) String() string {
	return fmt.Sprintf("{bucket: '%s', object: '%s'}", o.Bucket, o.Path)
}

func (o *ObjectT) StructHash() string {
	return hex.EncodeToString(md5.New().Sum([]byte(o.String())))
}
