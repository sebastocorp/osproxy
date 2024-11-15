package managers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"osproxy/api/v1alpha5"
	"osproxy/internal/objectstorage"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

type GCSManagerT struct {
	client   http.Client
	creds    *google.Credentials
	endpoint string
}

func (m *GCSManagerT) Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error) {
	m.endpoint = config.GCS.Endpoint
	if m.endpoint == "" {
		m.endpoint = "https://storage.googleapis.com"
	}
	endpointParts := strings.Split(m.endpoint, "://")
	if len(endpointParts) != 2 {
		err = fmt.Errorf("invalid endpoint format in gcs configuration, must be '<protocol>://<host>'")
		return err
	}
	if endpointParts[0] != "http" && endpointParts[0] != "https" {
		err = fmt.Errorf("invalid endpoint protocol in gcs configuration, must be 'http' or 'https'")
		return err
	}

	m.client.Timeout = 10 * time.Second
	credsBytes, err := base64.RawStdEncoding.DecodeString(config.GCS.Base64Credentials)
	if err != nil {
		return err
	}

	m.creds, err = google.CredentialsFromJSON(ctx, credsBytes, storage.CloudPlatformScope)

	return err
}

func (m *GCSManagerT) GetObject(obj objectstorage.ObjectT) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", m.endpoint, obj.Bucket, obj.Path), nil)
	if err != nil {
		return resp, err
	}

	token, err := m.creds.TokenSource.Token()
	if err != nil {
		return resp, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	resp, err = m.client.Do(req)

	return resp, err
}
