package osproxy

import (
	"fmt"
	"slices"
	"sync"

	"osproxy/api/v1alpha5"
	"osproxy/internal/osproxy/components/proxycomp"
)

type OSProxyT struct {
	config v1alpha5.OSProxyConfigT
	proxy  proxycomp.ProxyT
}

func NewOSProxy(configFilepath string) (o OSProxyT, err error) {
	o.config, err = parseConfig(configFilepath)
	if err != nil {
		return o, err
	}

	//--------------------------------------------------------------
	// Check config
	//--------------------------------------------------------------

	for _, src := range o.config.Proxy.Sources {
		srcTypes := []string{"S3", "GCS", "HTTP"}
		if !slices.Contains(srcTypes, src.Type) {
			err = fmt.Errorf("sources must be one of this types %v", srcTypes)
			return o, err
		}
	}

	for _, mod := range o.config.Proxy.RequestModifiers {
		modTypes := []string{"Path", "PathRegex", "Header"}
		if !slices.Contains(modTypes, mod.Type) {
			err = fmt.Errorf("modifiers must be one of this types %v", modTypes)
			return o, err
		}
	}

	routeTypes := []string{"Host", "PathPrefix", "HeaderValue"}
	if !slices.Contains(routeTypes, o.config.Proxy.RequestRouting.MatchType) {
		err = fmt.Errorf("config field 'proxy.requestRouting.matchType' must be one of this types %v", routeTypes)
		return o, err
	}

	if o.config.Proxy.RequestRouting.MatchType == "HeaderValue" && o.config.Proxy.RequestRouting.HeaderKey == "" {
		err = fmt.Errorf("header name in headerValue request routing match type must be set")
		return o, err
	}

	for _, route := range o.config.Proxy.RequestRouting.Routes {
		srcFound := false
		for _, sourcev := range o.config.Proxy.Sources {
			if sourcev.Name == route.Source {
				srcFound = true
				break
			}
		}
		if !srcFound {
			err = fmt.Errorf("unexisting source reference in routes config")
			return o, err
		}

		for _, routeModv := range route.Modifiers {
			modFound := false
			for _, modv := range o.config.Proxy.RequestModifiers {
				if modv.Name == routeModv {
					modFound = true
					break
				}
			}
			if !modFound {
				err = fmt.Errorf("unexisting modifier reference in routes config")
				return o, err
			}
		}
	}

	//--------------------------------------------------------------
	// Create components
	//--------------------------------------------------------------

	o.proxy, err = proxycomp.NewProxy(&o.config)
	if err != nil {
		return o, err
	}

	return o, err
}

func (o *OSProxyT) Run() (err error) {
	osWg := sync.WaitGroup{}
	osWg.Add(1)

	go o.proxy.Run(&osWg)

	osWg.Wait()

	return err
}
