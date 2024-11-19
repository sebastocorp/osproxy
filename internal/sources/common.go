package sources

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
)

type ObjectT struct {
	Bucket   string      `json:"bucket"`
	Path     string      `json:"path"`
	Metadata http.Header `json:"metadata"`
}

func (o *ObjectT) String() string {
	return fmt.Sprintf("{bucket: '%s', object: '%s'}", o.Bucket, o.Path)
}

func (o *ObjectT) StructHash() string {
	return hex.EncodeToString(md5.New().Sum([]byte(o.String())))
}
