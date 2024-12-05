package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

type RequestT struct {
	Method  string      `json:"method"`
	Host    string      `json:"host"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
}

type ResponseT struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Headers http.Header `json:"headers"`
	Request RequestT    `json:"request"`
}

func RequestID(r *http.Request) string {
	headers := "{"
	for hk, hvs := range r.Header {
		headers += "(" + hk + fmt.Sprintf("%v", hvs) + ")"

	}
	headers += "}"

	reqStr := fmt.Sprintf("{method: '%s', host: '%s', path: '%s', headers: '%s'}", r.Method, r.Host, r.URL.Path, headers)
	md5Hash := md5.New()
	_, err := md5Hash.Write([]byte(reqStr))
	if err != nil {
		return "UnableGetRequestID"
	}

	return hex.EncodeToString(md5Hash.Sum(nil))
}

func RequestStruct(r *http.Request) (req RequestT) {
	req.Method = r.Method
	req.Host = r.Host
	req.Path = r.URL.Path
	req.Headers = make(http.Header)
	for hk, hvs := range r.Header {
		for _, hv := range hvs {
			req.Headers.Add(hk, hv)
		}
	}

	return req
}

func ResponseStruct(r *http.Response) (res ResponseT) {
	res.Status = r.Status
	res.Code = r.StatusCode
	res.Headers = make(http.Header)
	for hk, hvs := range r.Header {
		for _, hv := range hvs {
			res.Headers.Add(hk, hv)
		}
	}

	res.Request = RequestStruct(r.Request)

	return res
}

func DefaultRequestStruct() (req RequestT) {
	req.Method = "none"
	req.Host = "none"
	req.Path = "none"
	req.Headers = make(http.Header)
	return req
}

func DefaultResponseStruct() (res ResponseT) {
	res.Code = 0
	res.Status = "none"
	res.Headers = make(http.Header)
	res.Request = DefaultRequestStruct()
	return res
}
