package backend

import (
	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/smtp"
	log "github.com/sirupsen/logrus"
)

// The package's init method
func init() {
	log.Debug("Init db storage backend package")
}

type DbStoreStorageBackend struct {
	config  *common.ApplicationConfig
	isReady bool
}

func (backend *DbStoreStorageBackend) IsStarted() bool {
	return backend.isReady
}

func (backend *DbStoreStorageBackend) Start() error {
	log.Info("Starting db backed mailstore backend instance...")

	return nil
}

func (backend *DbStoreStorageBackend) OnSubmitQueue(message *smtp.SmtpServerChannelMessage) {
	panic("implement me")
}
