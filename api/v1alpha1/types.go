package v1alpha1

type OSProxyConfigT struct {
	TransferService TransferServiceT `yaml:"transferService"`
	Relation        RelationT        `yaml:"relation"`
	OSConfig        OSConfigT        `yaml:"osConfig"`
}

type TransferServiceT struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Endpoint string `yaml:"endpoint"`
}

type RelationT struct {
	Type    string                       `yaml:"type"`
	Buckets map[string]FrontBackBucketsT `yaml:"buckets"`
}

type OSConfigT struct {
	S3  S3T  `yaml:"s3"`
	GCS GCST `yaml:"gcs"`
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

type S3T struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
}

type GCST struct {
	CredentialsFile string `yaml:"credentialsFile"`
}
