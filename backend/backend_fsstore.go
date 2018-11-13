package backend

import (
	"github.com/milhaux/common"
	log "github.com/sirupsen/logrus"
)

// The package's init method
func init() {
	log.Debug("Init fs storage backend package")
}

type FsStoreStorageBackend struct {
	config *common.ApplicationConfig
}

func (backend *FsStoreStorageBackend) Start() error {
	log.Info("Starting fs backed mailstore backend instance...")

	return nil
}
