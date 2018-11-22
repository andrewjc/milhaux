package backend

import (
	"github.com/andrewjc/milhaux/common"
	log "github.com/sirupsen/logrus"
)

// The package's init method
func init() {
	log.Debug("Init mem storage backend package")
}

type MemStoreStorageBackend struct {
	config *common.ApplicationConfig
}

func (backend *MemStoreStorageBackend) Start() error {
	log.Info("Starting memory backed mailstore backend instance...")

	return nil
}

func (backend *MemStoreStorageBackend) QueueSubmit(mailmessage *common.MailMessage) {

}
