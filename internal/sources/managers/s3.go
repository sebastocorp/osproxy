package managers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"osproxy/api/v1alpha5"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type S3ManagerT struct {
	creds  aws.Credentials
	signer *v4.Signer
	client *http.Client

	endpoint         string
	region           string
	emptyPayloadHash string
}

func (m *S3ManagerT) Init(ctx context.Context, config v1alpha5.ProxySourceConfigT) (err error) {
	m.emptyPayloadHash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	m.client = &http.Client{}
	m.client.Timeout = 10 * time.Second
	m.signer = v4.NewSigner()

	m.endpoint = config.S3.Endpoint
	endpointParts := strings.Split(m.endpoint, "://")
	if len(endpointParts) != 2 {
		err = fmt.Errorf("invalid endpoint format in s3 configuration, must be '<protocol>://<host>'")
		return err
	}
	if endpointParts[0] != "http" && endpointParts[0] != "https" {
		err = fmt.Errorf("invalid endpoint protocol in s3 configuration, must be 'http' or 'https'")
		return err
	}

	m.region = config.S3.Region
	credsProv := credentials.NewStaticCredentialsProvider(config.S3.AccessKeyID, config.S3.SecretAccessKey, "")
	m.creds, err = credsProv.Retrieve(context.Background())

	return err
}

func (m *S3ManagerT) GetObject(r *http.Request, bucket string) (resp *http.Response, err error) {
	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s/%s%s", m.endpoint, bucket, r.URL.Path), nil)
	if err != nil {
		return resp, err
	}

	req.Header.Set("x-amz-content-sha256", m.emptyPayloadHash)
	signingTime := time.Now().UTC()
	err = m.signer.SignHTTP(context.Background(), m.creds, req, m.emptyPayloadHash, "s3", m.region, signingTime)
	if err != nil {
		return resp, err
	}

	for hk, hvs := range r.Header {
		if _, ok := req.Header[hk]; !ok {
			for _, hv := range hvs {
				req.Header.Add(hk, hv)
			}
		}
	}

	resp, err = m.client.Do(req)

	return resp, err
}
