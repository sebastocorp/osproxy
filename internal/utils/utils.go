package utils

import (
	"fmt"
	"net/http"
)

const (
	DefaultSourceKey = "osproxy-source-default"
)

func RequestString(r *http.Request) string {
	headers := "{"
	for hk, hvs := range r.Header {
		headers += "(" + hk + fmt.Sprintf("%v", hvs) + ")"

	}
	headers += "}"
	return fmt.Sprintf("{method: '%s', url: '%s', headers: '%s'}", r.Method, r.URL.String(), headers)
}

func ResponseString(r *http.Response) string {
	headers := "{"
	for hk, hvs := range r.Header {
		headers += "(" + hk + fmt.Sprintf("%v", hvs) + ")"

	}
	headers += "}"
	return fmt.Sprintf("{status: '%s', code: '%d', headers: '%s'}", r.Status, r.StatusCode, headers)
}
