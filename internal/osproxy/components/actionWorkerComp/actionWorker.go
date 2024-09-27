package actionWorkerComp

import (
	"context"
	"log"
	"sync"
	"time"

	"osproxy/api/v1alpha3"
	"osproxy/internal/logger"
	"osproxy/internal/pools"
)

const (
	logExtraFieldKeyRequest       = "request"
	logExtraFieldKeyObject        = "object"
	logExtraFieldKeyBackendObject = "object_backend"
	logExtraFieldKeyError         = "error"
)

type ActionWorkerT struct {
	Config v1alpha3.ActionWorkerConfigT

	Log logger.LoggerT

	actionPool *pools.ActionPoolT
}

func NewActionWorker(config v1alpha3.ActionWorkerConfigT, actionPool *pools.ActionPoolT) (aw ActionWorkerT, err error) {
	aw.Config = config
	aw.actionPool = actionPool

	level, err := logger.GetLevel(aw.Config.Loglevel)
	if err != nil {
		log.Fatalf("unable to get log level in action worker config: %s", err.Error())
	}

	aw.Log = logger.NewLogger(context.Background(), level, map[string]any{
		"service":   "osproxy",
		"component": "actionWorker",
	})

	return aw, err
}

func (a *ActionWorkerT) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	logExtraFields := map[string]any{
		logExtraFieldKeyRequest: "none",
		logExtraFieldKeyObject:  "none",
		logExtraFieldKeyError:   "none",
	}

	for {

		pool := a.actionPool.Get()

		for k, v := range pool {
			logExtraFields[logExtraFieldKeyRequest] = v.Request.String()
			logExtraFields[logExtraFieldKeyObject] = v.Object.String()

			a.actionPool.Remove(k)

			backendObject, err := v.Request.GetObjectFromSource(a.Config.Source)
			if err != nil {
				logExtraFields[logExtraFieldKeyError] = err.Error()
				a.Log.Error("unable to get object from source config", logExtraFields)
			}
			logExtraFields[logExtraFieldKeyBackendObject] = backendObject.String()

			err = a.makeAPICall(v.Object, backendObject)
			if err != nil {
				logExtraFields[logExtraFieldKeyError] = err.Error()
				a.Log.Error("unable make api call", logExtraFields)
			}
		}

		time.Sleep(2 * time.Second)
	}

}
