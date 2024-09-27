package actionWorkerComp

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"osproxy/api/v1alpha3"
	"osproxy/internal/logger"
	"osproxy/internal/objectStorage"
	"osproxy/internal/pools"
)

const (
	logExtraFieldKeyObject     = "object"
	logExtraFieldKeyActionType = "action_type"
	logExtraFieldKeyError      = "error"

	actionTypeRequest = "request"
)

type ActionWorkerT struct {
	config v1alpha3.ActionWorkerConfigT
	log    logger.LoggerT

	actionPool  *pools.ActionPoolT
	actionFuncs map[string]func(objectStorage.ObjectT) error
}

func NewActionWorker(config v1alpha3.ActionWorkerConfigT, actionPool *pools.ActionPoolT) (aw ActionWorkerT, err error) {
	aw.config = config
	aw.actionPool = actionPool

	level, err := logger.GetLevel(aw.config.Loglevel)
	if err != nil {
		log.Fatalf("unable to get log level in action worker config: %s", err.Error())
	}

	aw.log = logger.NewLogger(context.Background(), level, map[string]any{
		"service":   "osproxy",
		"component": "actionWorker",
	})

	aw.actionFuncs = map[string]func(objectStorage.ObjectT) error{
		actionTypeRequest: aw.makeRequestAction,
	}

	if _, ok := aw.actionFuncs[aw.config.Type]; !ok {
		err = fmt.Errorf("action worker type not suported")
	}

	return aw, err
}

func (a *ActionWorkerT) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	logExtraFields := map[string]any{
		logExtraFieldKeyObject:     "none",
		logExtraFieldKeyError:      "none",
		logExtraFieldKeyActionType: a.config.Type,
	}

	for {
		pool := a.actionPool.Get()

		for k, v := range pool {
			logExtraFields[logExtraFieldKeyObject] = v.Object.String()

			a.actionPool.Remove(k)

			err := a.actionFuncs[a.config.Type](v.Object)
			if err != nil {
				logExtraFields[logExtraFieldKeyError] = err.Error()
				a.log.Error("unable make action", logExtraFields)
				continue
			}

			a.log.Debug("success in make action", logExtraFields)
		}

		time.Sleep(a.config.ScrapeInterval)
	}

}
