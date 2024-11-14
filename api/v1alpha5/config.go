package v1alpha5

type OSProxyConfigT struct {
	Proxy        ProxyConfigT        `yaml:"proxy"`
	ActionWorker ActionWorkerConfigT `yaml:"actionWorker"`
}
