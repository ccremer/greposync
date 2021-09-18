package valuestore

import (
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
)

type ValueStoreInstrumentation struct {
	log   logr.Logger
	scope string
}

func NewValueStoreInstrumentation(factory logging.LoggerFactory) *ValueStoreInstrumentation {
	return &ValueStoreInstrumentation{
		log: factory.NewGenericLogger(""),
	}
}

func (i *ValueStoreInstrumentation) attemptingLoadConfig(scope string, fileName string) {
	if i == nil {
		return
	}
	i.log.WithName(scope).V(logging.LevelDebug).Info("Loading config", "file", fileName)
}

func (i *ValueStoreInstrumentation) loadedConfigIfNil(scope string, fileName string, err error) error {
	if i == nil {
		return nil
	}
	if err != nil {
		i.log.WithName(scope).V(logging.LevelWarn).Info("file not loaded", "file", fileName, "error", err.Error())
	}
	return nil
}
