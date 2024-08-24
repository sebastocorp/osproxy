package utils

import (
	"fmt"
	"strings"
)

type RequestT struct {
	Host        string
	Port        string
	Path        string
	QueryParams map[string]string
}

func NewRequest(host, fullpath string) (r RequestT) {
	pathQueryParts := strings.SplitN(fullpath, "?", 2)

	r.Path = pathQueryParts[0]
	r.QueryParams = map[string]string{}

	if len(pathQueryParts) == 2 {
		queryParams := strings.Split(pathQueryParts[1], "&")
		for _, qp := range queryParams {
			qpParts := strings.SplitN(qp, "=", 2)
			if len(qpParts) == 2 {
				r.QueryParams[qpParts[0]] = qpParts[1]
			}
		}
	}

	hostParts := strings.Split(host, ":")

	r.Host = hostParts[0]
	if len(hostParts) == 2 {
		r.Port = hostParts[1]
	}

	return r
}

func (r *RequestT) String() string {
	qp := "{"
	for k, v := range r.QueryParams {
		qp += fmt.Sprintf("[%s:%s]", k, v)
	}
	qp += "}"
	return fmt.Sprintf("{host: '%s:%s', path: '%s', queryParams: '%s'}", r.Host, r.Port, r.Path, qp)
}
