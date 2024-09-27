package osproxy

import (
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

	actionPool := pools.NewActionPool()

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
