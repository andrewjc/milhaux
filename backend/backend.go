package backend

import (
	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/smtp"
	log "github.com/sirupsen/logrus"
)

var packageConfig = DefaultBackendConfig{}

type DefaultBackendConfig struct {
	backend   int8
	available bool
}

const (
	MEMSTORE = iota
	FILESTORE
	DBSTORE
)

// The package's init method
func init() {
	log.Debug("Init mailstore package")

	// Override config values for testing.
	packageConfig.backend = FILESTORE
	packageConfig.available = true
}

type MailStoreBackend struct {
	smtpComponentChannel chan common.MailMessage
}

func NewMailStoreBackend(config *common.ApplicationConfig) MailStoreBackend {
	log.Debug("Creating a mailstore backend instance...")

	return MailStoreBackend{}
}

func (s *MailStoreBackend) Start() error {
	// Setup a listener for the communication channel with
	return nil
}

func (s *MailStoreBackend) OnSubmitQueue(message *smtp.SmtpServerChannelMessage) {
	log.Debug("Got message to backend...")

	message.Data.QueueId = "123123"
}
