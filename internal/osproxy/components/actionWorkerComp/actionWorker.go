package actionWorkerComp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"osproxy/api/v1alpha5"
	"osproxy/internal/global"
	"osproxy/internal/logger"
	"osproxy/internal/objectstorage"
	"osproxy/internal/pools"
)

const (
	actionTypeRequest = "request"
)

type ActionWorkerT struct {
	config *v1alpha5.OSProxyConfigT
	log    logger.LoggerT

	actionPool  *pools.ActionPoolT
	actionFuncs map[string]func(objectstorage.ObjectT) error
}

func NewActionWorker(config *v1alpha5.OSProxyConfigT, actionPool *pools.ActionPoolT) (aw ActionWorkerT, err error) {
	aw.config = config
	aw.actionPool = actionPool

	logCommon := global.GetLogCommonFields()
	logCommon[global.LogFieldKeyCommonComponent] = global.LogFieldValueComponentActionWorker
	aw.log = logger.NewLogger(context.Background(), logger.GetLevel(aw.config.ActionWorker.Loglevel), logCommon)

	aw.actionFuncs = map[string]func(objectstorage.ObjectT) error{
		actionTypeRequest: aw.makeRequestAction,
	}

	if _, ok := aw.actionFuncs[aw.config.ActionWorker.Type]; !ok {
		err = fmt.Errorf("action worker type not suported")
	}

	return aw, err
}

func (a *ActionWorkerT) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	logExtraFields := global.GetLogExtraFieldsActionWorker()
	logExtraFields[global.LogFieldKeyExtraActionType] = a.config.ActionWorker.Type

	emptyPoolLog := true
	for {
		logExtraFields[global.LogFieldKeyExtraError] = global.LogFieldValueDefault
		logExtraFields[global.LogFieldKeyExtraObject] = global.LogFieldValueDefault

		pool := a.actionPool.Get()

		if len(pool) == 0 {
			if emptyPoolLog {
				a.log.Debug("wait for empty pool", logExtraFields)
			}
			emptyPoolLog = false
			time.Sleep(a.config.ActionWorker.ScrapeInterval)
			continue
		}
		emptyPoolLog = true

		for key, request := range pool {
			a.actionPool.Remove(key)

			logExtraFields[global.LogFieldKeyExtraError] = global.LogFieldValueDefault
			logExtraFields[global.LogFieldKeyExtraObject] = request.Object.String()
			err := a.actionFuncs[a.config.ActionWorker.Type](request.Object)
			if err != nil {
				logExtraFields[global.LogFieldKeyExtraError] = err.Error()
				a.log.Error("unable make action", logExtraFields)
				continue
			}

			a.log.Info("success in make action", logExtraFields)
		}

		time.Sleep(a.config.ActionWorker.ScrapeInterval)
	}

}
