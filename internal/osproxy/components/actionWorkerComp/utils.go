package actionWorkerComp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"osproxy/internal/objectstorage"
)

func (a *ActionWorkerT) makeRequestAction(Object objectstorage.ObjectT) (err error) {
	bodyBytes, err := json.Marshal(Object)
	if err != nil {
		return err
	}

	http.DefaultClient.Timeout = 200 * time.Millisecond
	resp, err := http.Post(a.config.ActionWorker.Request.URL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
