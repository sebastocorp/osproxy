package utils

import (
	"fmt"
	"net/http"
)

const (
	DefaultSourceKey = "osproxy-source-default"
)

func RequestString(r *http.Request) string {
	return fmt.Sprintf("{url: '%s', method: '%s'}", r.URL.String(), r.Method)
}
