package objectstorage

import (
	"context"
	"encoding/hex"
	"osproxy/api/v1alpha5"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSManagerT struct {
	ctx    context.Context
	client *storage.Client
}

type GCSReaderT struct {
	reader *storage.Reader
}

func (m *GCSManagerT) Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error) {
	m.ctx = ctx
	m.client, err = storage.NewClient(m.ctx, option.WithCredentialsFile(config.GCS.CredentialsFile))
	return err
}

func (m *GCSManagerT) GetObject(obj ObjectT) (objReader ObjectReaderI, info ObjectInfoT, err error) {
	gcsobj := m.client.Bucket(obj.Bucket).Object(obj.Path)
	stat, err := gcsobj.Attrs(m.ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			info.NotExistError = true
		}
		return objReader, info, err
	}

	info.MD5 = hex.EncodeToString(stat.MD5)
	info.Size = stat.Size
	info.ContentType = stat.ContentType

	gcsReader, err := gcsobj.NewReader(m.ctx)
	objReader = &GCSReaderT{
		reader: gcsReader,
	}

	return objReader, info, err
}

func (r *GCSReaderT) Close() error {
	return r.reader.Close()
}

func (r *GCSReaderT) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	return n, err
}
