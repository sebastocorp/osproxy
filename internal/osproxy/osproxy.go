package osproxy

import (
	"fmt"
	"slices"
	"sync"

	"osproxy/api/v1alpha5"
	"osproxy/internal/osproxy/components/actionWorkerComp"
	"osproxy/internal/osproxy/components/proxyComp"
	"osproxy/internal/pools"
)

type OSProxyT struct {
	config       v1alpha5.OSProxyConfigT
	proxy        proxyComp.ProxyT
	actionWorker actionWorkerComp.ActionWorkerT
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
		if !slices.Contains([]string{"s3", "gcs"}, src.Type) {
			err = fmt.Errorf("sources must be one of this types %v", []string{"s3", "gcs"})
			return o, err
		}
	}

	for _, mod := range o.config.Proxy.RequestModifiers {
		if !slices.Contains([]string{"path", "header"}, mod.Type) {
			err = fmt.Errorf("modifiers must be one of this types %v", []string{"path", "header"})
			return o, err
		}
	}

	if !slices.Contains([]string{"host", "pathPrefix", "headerValue"}, o.config.Proxy.RequestRouting.MatchType) {
		err = fmt.Errorf("config field 'proxy.requestRouting.matchType' must be one of this types %v", []string{"host", "pathPrefix", "headerValue"})
		return o, err
	}

	if o.config.Proxy.RequestRouting.MatchType == "headerValue" && o.config.Proxy.RequestRouting.HeaderKey == "" {
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
	actionPool := pools.NewActionPool(o.config.ActionWorker.PoolCapacity)

	o.proxy, err = proxyComp.NewProxy(&o.config, actionPool)
	if err != nil {
		return o, err
	}

	o.actionWorker, err = actionWorkerComp.NewActionWorker(&o.config, actionPool)
	if err != nil {
		return o, err
	}

	return o, err
}

func (o *OSProxyT) Run() (err error) {
	osWg := sync.WaitGroup{}
	osWg.Add(2)

	go o.proxy.Run(&osWg)
	go o.actionWorker.Run(&osWg)

	osWg.Wait()

	return err
}
