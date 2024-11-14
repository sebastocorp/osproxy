package objectstorage

import (
	"context"
	"osproxy/api/v1alpha5"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3ManagerT struct {
	Ctx    context.Context
	Client *minio.Client
}

type S3ReaderT struct {
	reader *minio.Object
}

func (m *S3ManagerT) Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error) {
	m.Ctx = ctx
	m.Client, err = minio.New(
		config.S3.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(config.S3.AccessKeyID, config.S3.SecretAccessKey, ""),
			Region: config.S3.Region,
			Secure: config.S3.Secure,
		},
	)

	return err
}

func (m *S3ManagerT) GetObject(obj ObjectT) (objReader ObjectReaderI, info ObjectInfoT, err error) {
	object, err := m.Client.GetObject(context.Background(), obj.Bucket, obj.Path, minio.GetObjectOptions{})
	if err != nil {
		return objReader, info, err
	}

	stat, err := object.Stat()
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			info.NotExistError = true
		}
		return objReader, info, err
	}

	info.MD5 = stat.ETag
	info.ContentType = stat.ContentType
	info.Size = stat.Size

	objReader = &S3ReaderT{
		reader: object,
	}

	return objReader, info, err
}

func (r *S3ReaderT) Close() error {
	return r.reader.Close()
}

func (r *S3ReaderT) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	return n, err
}
