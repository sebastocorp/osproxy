package osproxy

import (
	"fmt"
	"sync"

	"osproxy/internal/osproxy/components/actionWorkerComp"
	"osproxy/internal/osproxy/components/proxyComp"
	"osproxy/internal/pools"
)

type OSProxyT struct {
	proxy        proxyComp.ProxyT
	actionWorker actionWorkerComp.ActionWorkerT
}

func NewOSProxy(configFilepath string) (o OSProxyT, err error) {
	osConfig, err := parseConfig(configFilepath)
	if err != nil {
		return o, err
	}

	//--------------------------------------------------------------
	// Check config
	//--------------------------------------------------------------
	if osConfig.Proxy.RequestRouting.Type == "headerValue" && osConfig.Proxy.RequestRouting.HeaderName == "" {
		err = fmt.Errorf("header name in headerValue request routing type must be set")
		return o, err
	}

	//--------------------------------------------------------------
	// Create components
	//--------------------------------------------------------------
	actionPool := pools.NewActionPool(osConfig.ActionWorker.PoolCapacity)

	o.proxy, err = proxyComp.NewProxy(osConfig.Proxy, actionPool)
	if err != nil {
		return o, err
	}

	o.actionWorker, err = actionWorkerComp.NewActionWorker(osConfig.ActionWorker, actionPool)
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
