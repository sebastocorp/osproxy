package objectStorage

import (
	"context"
	"fmt"
	"reflect"

	"cloud.google.com/go/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/api/option"
)

type ManagerT struct {
	Ctx context.Context
	S3  S3T
	GCS GCST
}

type S3T struct {
	Client          *minio.Client
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

type GCST struct {
	Client          *storage.Client
	CredentialsFile string
}

type ObjectT struct {
	BucketName string `json:"bucket"`
	ObjectPath string `json:"path"`
	Etag       string `json:"etag"`
}

type ObjectInfoT struct {
	Size        int64
	ContentType string
}

func NewManager(ctx context.Context, s3 S3T, gcs GCST) (man ManagerT, err error) {
	s3ConfigEmpty := reflect.ValueOf(s3).IsZero()
	gcsConfigEmpty := reflect.ValueOf(gcs).IsZero()
	if s3ConfigEmpty && gcsConfigEmpty {
		err = fmt.Errorf("both s3 and gcs config are empty")
		return man, err
	}

	man.Ctx = ctx

	if !s3ConfigEmpty {
		man.S3 = s3
		man.S3.Client, err = minio.New(
			man.S3.Endpoint,
			&minio.Options{
				Creds:  credentials.NewStaticV4(man.S3.AccessKeyID, man.S3.SecretAccessKey, ""),
				Secure: true,
			},
		)
		if err != nil {
			return man, err
		}
	}

	if !gcsConfigEmpty {
		man.GCS.CredentialsFile = gcs.CredentialsFile
		man.GCS.Client, err = storage.NewClient(man.Ctx, option.WithCredentialsFile(man.GCS.CredentialsFile))
	}

	return man, err
}

func (m *ManagerT) GCSObjectExist(obj ObjectT) (exist bool, info ObjectInfoT, err error) {
	exist = true
	stat, err := m.GCS.Client.Bucket(obj.BucketName).Object(obj.ObjectPath).Attrs(m.Ctx)
	if err != nil {
		if err != storage.ErrObjectNotExist {
			return exist, info, err
		}
		err = nil
		exist = false
	}

	if exist {
		info = ObjectInfoT{
			Size:        stat.Size,
			ContentType: stat.ContentType,
		}
	}

	return exist, info, err
}

func (m *ManagerT) S3ObjectExist(obj ObjectT) (exist bool, info ObjectInfoT, err error) {
	exist = true
	stat, err := m.S3.Client.StatObject(m.Ctx, obj.BucketName, obj.ObjectPath, minio.GetObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code != "NoSuchKey" {
			return exist, info, err
		}
		err = nil
		exist = false
	}

	if exist {
		info = ObjectInfoT{
			Size:        stat.Size,
			ContentType: stat.ContentType,
		}
	}

	return exist, info, err
}

func (o *ObjectT) String() string {
	return fmt.Sprintf("{bucket: '%s', object: '%s'}", o.BucketName, o.ObjectPath)
}
