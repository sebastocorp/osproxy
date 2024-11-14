package v1alpha5

//--------------------------------------------------------------
// PROXY CONFIG
//--------------------------------------------------------------

type ProxyConfigT struct {
	Loglevel       string                              `yaml:"loglevel"`
	Address        string                              `yaml:"address"`
	Port           string                              `yaml:"port"`
	Modifications  map[string]ProxyModificationConfigT `yaml:"modifications"`
	Sources        map[string]ProxySourceConfigT       `yaml:"sources"`
	RequestRouting ProxyRequestRoutingConfigT          `yaml:"requestRouting"`
}

// Modifications config

type ProxyModificationConfigT struct {
	Type   string                         `yaml:"type"`
	Path   ProxyModificationPathConfigT   `yaml:"path"`
	Header ProxyModificationHeaderConfigT `yaml:"header"`
}

type ProxyModificationPathConfigT struct {
	AddPrefix    string `yaml:"addPrefix"`
	RemovePrefix string `yaml:"removePrefix"`
}

type ProxyModificationHeaderConfigT struct {
	Name   string `yaml:"name"`
	Remove bool   `yaml:"remove"`
	Value  string `yaml:"value"`
}

// Sources config

type ProxySourceConfigT struct {
	Type string                `yaml:"type"`
	S3   ProxySourceS3ConfigT  `yaml:"s3,omitempty"`
	GCS  ProxySourceGCSConfigT `yaml:"gcs,omitempty"`
}

type ProxySourceS3ConfigT struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Region          string `yaml:"region"`
	Secure          bool   `yaml:"secure"`
}

type ProxySourceGCSConfigT struct {
	CredentialsFile string `yaml:"credentialsFile"`
}

// Routing config

type ProxyRequestRoutingConfigT struct {
	MatchType string                       `yaml:"matchType"`
	HeaderKey string                       `yaml:"headerKey"`
	Routes    map[string]ProxyRouteConfigT `yaml:"routes"`
}

type ProxyRouteConfigT struct {
	Modifications []string `yaml:"modifications"`
	Source        string   `yaml:"source"`
	Bucket        string   `yaml:"bucket"`
}
