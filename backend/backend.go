package backend

import (
	"github.com/andrewjc/milhaux/common"
	"github.com/andrewjc/milhaux/smtp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	MEMSTORE = iota
	FILESTORE
	DBSTORE
)

// The package's init method
func init() {
	log.Debug("Init mailstore package")

}

type Mailbox interface {
	Name() string
}

type MailStoreStorageProvider interface {
	Start() error
	IsStarted() bool
	Store(message *common.MailMessage) error
	Mailboxes() []Mailbox
}

type MailStoreBackend struct {
	storageConnector common.StorageConnector
	commandProcessor *BackendCommandProcessor

	storageProvider MailStoreStorageProvider

	messageQueue chan smtp.SmtpServerChannelMessage
	workerPool   chan chan smtp.SmtpServerChannelMessage
	maxWorkers   int
}

func NewMailStoreBackend(config *common.ApplicationConfig) *MailStoreBackend {
	log.Debug("Creating a mailstore backend instance...")

	backend := &MailStoreBackend{}
	backend.maxWorkers = config.GetSmtpServerConfig().SMTP_OPTION_MAX_QUEUE_WORKERS
	backend.workerPool = make(chan chan smtp.SmtpServerChannelMessage, config.GetSmtpServerConfig().SMTP_OPTION_MAX_QUEUE_WORKERS)

	backend.storageProvider = &MemStoreStorageBackend{}

	backend.commandProcessor = &BackendCommandProcessor{}

	// If we're compiled in monolithic mode, use an in-process storage connector that assumes
	// that all components are running in the same process space.
	if config.GetBackendConfig().ListenInterface == "embedded" {
		backend.storageConnector = &InMemoryStorageConnector{*backend}
	}

	/*switch {
	case packageConfig.backend == MEMSTORE:
		backend.storageProvider = &MemStoreStorageBackend{}
	case packageConfig.backend == FILESTORE:
		backend.storageProvider = &FsStoreStorageBackend{}
	case packageConfig.backend == DBSTORE:
		backend.storageProvider = &DbStoreStorageBackend{}
	}*/

	return backend
}

type QueueWorker struct {
	workerId       string
	backend        *MailStoreBackend
	messageChannel chan smtp.SmtpServerChannelMessage
	quit           chan bool
}

func (backend *MailStoreBackend) InitSmtpMessageChannelListener(channel chan smtp.SmtpServerChannelMessage) {
	backend.messageQueue = channel
}

func (backend *MailStoreBackend) Start() error {

	if err := backend.storageProvider.Start(); err != nil {
		log.Errorf("Error occurred while starting backend storage provider: %v", err.Error())
		return err
	}

	// starting n number of workers
	for i := 0; i < backend.maxWorkers; i++ {
		worker := NewQueueWorker(backend)
		worker.Start()
	}

	go backend.dispatch()

	return nil
}

func (backend *MailStoreBackend) dispatch() {
	for {
		select {
		case job := <-backend.messageQueue:
			// a job request has been received
			log.Debugf("Dispatcher got new job in request queue...")

			go func(job smtp.SmtpServerChannelMessage) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				log.Debugf("Dispatcher obtaining worker job channel...")
				jobChannel := <-backend.workerPool

				// dispatch the job to the worker job channel
				log.Debugf("Dispatcher job to worker...")
				jobChannel <- job
			}(job)
		}
	}
}

func (backend *MailStoreBackend) GetStorageConnector() common.StorageConnector {
	return backend.storageConnector
}

func NewQueueWorker(backend *MailStoreBackend) *QueueWorker {
	return &QueueWorker{
		workerId:       uuid.New().String(),
		backend:        backend,
		messageChannel: make(chan smtp.SmtpServerChannelMessage),
		quit:           make(chan bool)}
}

func (w *QueueWorker) Start() {
	log.Debugf("Starting queue worker [%v]...", w.workerId)
	go func() {
		for {
			// register the current worker into the worker queue.
			w.backend.workerPool <- w.messageChannel

			select {
			case message := <-w.messageChannel:
				// we have received a work request.
				w.onProcessQueueMessage(message)

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *QueueWorker) StopWorker() {
	go func() {
		w.quit <- true
	}()
}

func (w *QueueWorker) onProcessQueueMessage(message smtp.SmtpServerChannelMessage) {
	log.Debugf("Worker [%v] processing queue message %v", w.workerId, message.Data.QueueId)
	backend := w.backend

	backend.storageProvider.Store(message.Data)
}
