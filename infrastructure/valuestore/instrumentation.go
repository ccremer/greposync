package valuestore

import (
	"github.com/ccremer/greposync/infrastructure/logging"
	"github.com/go-logr/logr"
)

type ValueStoreInstrumentation struct {
	log logr.Logger
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
	i.log.WithName(scope).V(1).Info("Loading config", "file", fileName)
}

func (i *ValueStoreInstrumentation) loadedConfigIfNil(scope string, err error) error {
	if i == nil {
		return nil
	}
	if err != nil {
		i.log.WithName(scope).V(1).Info("File not loaded", "error", err.Error())
	}
	return nil
}
