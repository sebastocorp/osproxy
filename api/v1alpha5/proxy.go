package v1alpha5

//--------------------------------------------------------------
// PROXY CONFIG
//--------------------------------------------------------------

type ProxyConfigT struct {
	Loglevel         string                     `yaml:"loglevel"`
	Protocol         string                     `yaml:"protocol"`
	Address          string                     `yaml:"address"`
	Port             string                     `yaml:"port"`
	Sources          []ProxySourceConfigT       `yaml:"sources"`
	RequestModifiers []ProxyModifierConfigT     `yaml:"requestModifiers"`
	RequestRouting   ProxyRequestRoutingConfigT `yaml:"requestRouting"`
	RespReactions    ProxyRespReactionsConfigT  `yaml:"responseReactions"`
}

// Sources config

type ProxySourceConfigT struct {
	Name string                `yaml:"name"`
	Type string                `yaml:"type"`
	S3   ProxySourceS3ConfigT  `yaml:"s3,omitempty"`
	GCS  ProxySourceGCSConfigT `yaml:"gcs,omitempty"`
}

type ProxySourceS3ConfigT struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Region          string `yaml:"region"`
}

type ProxySourceGCSConfigT struct {
	Endpoint          string `yaml:"endpoint"`
	Base64Credentials string `yaml:"base64Credentials"`
}

// Modifications config

type ProxyModifierConfigT struct {
	Name   string                     `yaml:"name"`
	Type   string                     `yaml:"type"`
	Path   ProxyModifierPathConfigT   `yaml:"path"`
	Header ProxyModifierHeaderConfigT `yaml:"header"`
}

type ProxyModifierPathConfigT struct {
	AddPrefix    string `yaml:"addPrefix"`
	RemovePrefix string `yaml:"removePrefix"`
}

type ProxyModifierHeaderConfigT struct {
	Name   string `yaml:"name"`
	Remove bool   `yaml:"remove"`
	Value  string `yaml:"value"`
}

// Routing config

type ProxyRequestRoutingConfigT struct {
	MatchType string                       `yaml:"matchType"`
	HeaderKey string                       `yaml:"headerKey"`
	Routes    map[string]ProxyRouteConfigT `yaml:"routes"`
}

type ProxyRouteConfigT struct {
	Source    string   `yaml:"source"`
	Modifiers []string `yaml:"modifiers"`
	Bucket    string   `yaml:"bucket"`
}

// Response reactions

type ProxyRespReactionsConfigT struct {
	Conditions []ProxyConditionConfigT `yaml:"conditions"`
	Reactions  []ProxyReactionConfigT  `yaml:"reactions"`
}

type ProxyConditionConfigT struct {
	Name string `yaml:"name"`
}

type ProxyReactionConfigT struct {
	Name string `yaml:"name"`
}
