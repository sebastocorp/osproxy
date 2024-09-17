package v1alpha2

type OSProxyConfigT struct {
	Proxy         ProxyT         `yaml:"proxy"`
	Action        ActionT        `yaml:"action"`
	ObjectStorage ObjectStorageT `yaml:"objectStorage"`
}

//--------------------------------------------------------------
// PROXY CONFIG
//--------------------------------------------------------------

type ProxyT struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}

//--------------------------------------------------------------
// ACTION CONFIG
//--------------------------------------------------------------

type ActionT struct {
	StatusCode int      `yaml:"statusCode"`
	APICall    APICallT `yaml:"apiCall"`
}

type APICallT struct {
	URL string `yaml:"url"`
}

//--------------------------------------------------------------
// OBJECT STORAGE CONFIG
//--------------------------------------------------------------

type ObjectStorageT struct {
	S3       S3T       `yaml:"s3,omitempty"`
	GCS      GCST      `yaml:"gcs,omitempty"`
	Relation RelationT `yaml:"relation"`
}

type S3T struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Region          string `yaml:"region"`
	Secure          bool   `yaml:"secure"`
}

type GCST struct {
	CredentialsFile string `yaml:"credentialsFile"`
}

//--------------------------------------------------------------
// RELATION CONFIG
//--------------------------------------------------------------

type RelationT struct {
	Type    string                       `yaml:"type"`
	Buckets map[string]FrontBackBucketsT `yaml:"buckets"`
}

type FrontBackBucketsT struct {
	Frontend BucketSubpathT `yaml:"frontend"`
	Backend  BucketSubpathT `yaml:"backend"`
}

type BucketSubpathT struct {
	BucketName       string `yaml:"bucketName"`
	AddPathPrefix    string `yaml:"addPathPrefix"`
	RemovePathPrefix string `yaml:"removePathPrefix"`
}
