package backend

import (
	"github.com/milhaux/common"
	log "github.com/sirupsen/logrus"
)

// The package's init method
func init() {
	log.Debug("Init db storage backend package")
}

type DbStoreStorageBackend struct {
	config *common.ApplicationConfig
}

func (backend *DbStoreStorageBackend) Start() error {
	log.Info("Starting db backed mailstore backend instance...")

	return nil
}
