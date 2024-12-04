package managers

import (
	"context"
	"fmt"
	"net/http"
	"osproxy/api/v1alpha5"
	"strings"
	"time"
)

type HTTPManagerT struct {
	client   *http.Client
	endpoint string
}

func (m *HTTPManagerT) Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error) {
	m.endpoint = config.HTTP.Endpoint
	endpointParts := strings.Split(m.endpoint, "://")
	if len(endpointParts) != 2 {
		err = fmt.Errorf("invalid endpoint format in http configuration, must be '<protocol>://<host>'")
		return err
	}
	if endpointParts[0] != "http" && endpointParts[0] != "https" {
		err = fmt.Errorf("invalid endpoint protocol in http configuration, must be 'http' or 'https'")
		return err
	}

	m.client = &http.Client{}
	m.client.Timeout = 10 * time.Second

	return err
}

func (m *HTTPManagerT) GetObject(r *http.Request, bucket string) (resp *http.Response, err error) {
	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s/%s%s", m.endpoint, bucket, r.URL.Path), r.Body)
	if err != nil {
		return resp, err
	}

	for hk, hvs := range r.Header {
		for _, hv := range hvs {
			req.Header.Add(hk, hv)
		}
	}

	resp, err = m.client.Do(req)

	return resp, err
}
