package backend

import (
	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/smtp"
	log "github.com/sirupsen/logrus"
)

// The package's init method
func init() {
	log.Debug("Init fs storage backend package")
}

type FsStoreStorageBackend struct {
	config  *common.ApplicationConfig
	isReady bool
}

func (backend *FsStoreStorageBackend) IsStarted() bool {
	return backend.isReady
}

func (backend *FsStoreStorageBackend) Start() error {
	log.Info("Starting fs backed mailstore backend instance...")

	return nil
}

func (backend *FsStoreStorageBackend) OnSubmitQueue(message *smtp.SmtpServerChannelMessage) {
	panic("implement me")
}
