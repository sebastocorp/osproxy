package proxycomp

import (
	"fmt"
	"net/http"
	"osproxy/api/v1alpha5"
	"osproxy/internal/utils"
	"strconv"
	"strings"
)

func (p *ProxyT) requestResponseError(respWriter http.ResponseWriter, respStatusCode int) string {
	respMessage := fmt.Sprintf("%d %s", respStatusCode, http.StatusText(respStatusCode))

	respError := &http.Response{
		Header:     http.Header{},
		StatusCode: respStatusCode,
		Status:     respMessage,
	}

	respError.Header.Set("Content-Type", "text/plain")
	respError.Header.Set("Content-Length", strconv.Itoa(len(respMessage)))

	// response to user request
	for hk, hvs := range respError.Header {
		for _, hv := range hvs {
			respWriter.Header().Set(hk, hv)
		}
	}
	respWriter.WriteHeader(respStatusCode)
	respWriter.Write([]byte(respError.Status))

	return utils.ResponseString(respError)
}

func (p *ProxyT) getRouteFromRequest(r *http.Request) (route v1alpha5.ProxyRouteConfigT, err error) {
	var found bool = false
	switch p.config.Proxy.RequestRouting.MatchType {
	case "host":
		{
			route, found = p.config.Proxy.RequestRouting.Routes[r.Host]
		}
	case "headerValue":
		{
			route, found = p.config.Proxy.RequestRouting.Routes[r.Header.Get(p.config.Proxy.RequestRouting.HeaderKey)]

		}
	case "pathPrefix":
		{
			requestPath := strings.SplitN(r.URL.Path, "?", 2)[0]
			for prefix, rout := range p.config.Proxy.RequestRouting.Routes {
				if strings.HasPrefix(requestPath, prefix) {
					route = rout
					found = true
					break
				}
			}
		}
	}

	if !found {
		err = fmt.Errorf("routing config not provided for this request")
		return route, err
	}

	return route, err
}

func (p *ProxyT) modRequest(r *http.Request, modifications []string) (err error) {
	r.URL.Path = strings.SplitN(r.URL.Path, "?", 2)[0]
	for _, modn := range modifications {
		mod := p.requestModifiers[modn]
		switch mod.Type {
		case "path":
			{
				r.URL.Path = mod.Path.AddPrefix + strings.TrimPrefix(r.URL.Path, mod.Path.RemovePrefix)
			}
		case "header":
			{
				r.Header.Set(mod.Header.Name, mod.Header.Value)
				if mod.Header.Remove {
					r.Header.Del(mod.Header.Name)
				}
			}
		}
	}

	return err
}
