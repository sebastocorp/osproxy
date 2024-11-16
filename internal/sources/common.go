package sources

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type ObjectT struct {
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
}

func (o *ObjectT) String() string {
	return fmt.Sprintf("{bucket: '%s', object: '%s'}", o.Bucket, o.Path)
}

func (o *ObjectT) StructHash() string {
	return hex.EncodeToString(md5.New().Sum([]byte(o.String())))
}
