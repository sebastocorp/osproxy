package actionWorkerComp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"osproxy/internal/objectStorage"
)

func (a *ActionWorkerT) makeAPICall(fObject, bObject objectStorage.ObjectT) (err error) {
	type TransferT struct {
		From objectStorage.ObjectT `json:"from"`
		To   objectStorage.ObjectT `json:"to"`
	}

	transfer := TransferT{
		From: bObject,
		To:   fObject,
	}

	bodyBytes, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	http.DefaultClient.Timeout = 200 * time.Millisecond
	resp, err := http.Post(a.Config.APICall.URL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
