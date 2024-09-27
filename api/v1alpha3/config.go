package v1alpha3

import "time"

type OSProxyConfigT struct {
	Proxy        ProxyConfigT        `yaml:"proxy"`
	ActionWorker ActionWorkerConfigT `yaml:"actionWorker"`
}

//--------------------------------------------------------------
// PROXY CONFIG
//--------------------------------------------------------------

type ProxyConfigT struct {
	Loglevel string        `yaml:"loglevel"`
	Address  string        `yaml:"address"`
	Port     string        `yaml:"port"`
	Source   SourceConfigT `yaml:"source"`
}

//--------------------------------------------------------------
// ACTION WORKER CONFIG
//--------------------------------------------------------------

type ActionWorkerConfigT struct {
	Loglevel       string         `yaml:"loglevel"`
	Type           string         `yaml:"type"`
	ScrapeInterval time.Duration  `yaml:"scrapeInterval"`
	Request        RequestConfigT `yaml:"request"`
}

type RequestConfigT struct {
	URL string `yaml:"url"`
}

//--------------------------------------------------------------
// SOURCE STORAGE CONFIG
//--------------------------------------------------------------

type SourceConfigT struct {
	Config  ObjectStorageConfigT           `yaml:"config"`
	Type    string                         `yaml:"type"`
	Buckets map[string]BucketObjectConfigT `yaml:"buckets"`
}

type BucketObjectConfigT struct {
	Bucket    string                    `yaml:"bucket"`
	ObjectMod ObjectModificationConfigT `yaml:"objectModification"`
}

type ObjectModificationConfigT struct {
	AddPrefix    string `yaml:"addPrefix"`
	RemovePrefix string `yaml:"removePrefix"`
}

//--------------------------------------------------------------
// OBJECT STORAGE CONFIG
//--------------------------------------------------------------

type ObjectStorageConfigT struct {
	S3  S3ConfigT  `yaml:"s3,omitempty"`
	GCS GCSConfigT `yaml:"gcs,omitempty"`
}

type S3ConfigT struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	Region          string `yaml:"region"`
	Secure          bool   `yaml:"secure"`
}

type GCSConfigT struct {
	CredentialsFile string `yaml:"credentialsFile"`
}
