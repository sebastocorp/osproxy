package proxyComp

import (
	"fmt"
	"net/http"
	"strconv"
)

func (p *ProxyT) requestResponseError(respWriter http.ResponseWriter, respStatusCode int, respMessage string) {
	respMessage = fmt.Sprintf("%d %s\n", respStatusCode, respMessage)

	// response to user request
	respWriter.Header().Set("Content-Type", "text/plain")
	respWriter.Header().Set("Content-Length", strconv.Itoa(len(respMessage)))
	respWriter.WriteHeader(respStatusCode)
	respWriter.Write([]byte(respMessage))
}
