package objectStorage

import (
	"context"
	"fmt"
	"reflect"

	"osproxy/api/v1alpha2"

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
	BucketName string      `json:"bucket"`
	ObjectPath string      `json:"path"`
	Info       ObjectInfoT `json:"-"`
}

type ObjectInfoT struct {
	Exist       bool
	MD5         string
	Size        int64
	ContentType string
}

func NewManager(ctx context.Context, s3 v1alpha2.S3T, gcs v1alpha2.GCST) (man ManagerT, err error) {
	man.Ctx = ctx

	man.S3Client, err = minio.New(
		s3.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(s3.AccessKeyID, s3.SecretAccessKey, ""),
			Region: s3.Region,
			Secure: s3.Secure,
		},
	)
	if err != nil {
		return man, err
	}

	if reflect.ValueOf(s3).IsZero() {
		man.GCSClient, err = storage.NewClient(man.Ctx, option.WithCredentialsFile(gcs.CredentialsFile))
	}

	return man, err
}

func (m *ManagerT) S3GetObject(obj ObjectT) (object *minio.Object, info ObjectInfoT, err error) {
	object, err = m.S3Client.GetObject(context.Background(), obj.BucketName, obj.ObjectPath, minio.GetObjectOptions{})
	if err != nil {
		return object, info, err
	}

	stat, err := object.Stat()
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			err = nil
			info.Exist = false
		}
		return object, info, err
	}

	info.Exist = true
	info.MD5 = stat.ETag
	info.ContentType = stat.ContentType
	info.Size = stat.Size

	return object, info, err
}

func (o *ObjectT) String() string {
	return fmt.Sprintf("{bucket: '%s', object: '%s'}", o.BucketName, o.ObjectPath)
}
