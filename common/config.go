package common

var appConfigInstance = CreateDefaultAppConfig()

func init() {
}

type ApplicationConfig struct {
	smtpServerConfig  *SmtpServerConfig
	imap4ServerConfig *Imap4ServerConfig
	backendConfig     *BackendConfig
	version           string
	loglevel          uint32
}

type SmtpServerConfig struct {
	Port                                    uint32
	Hostname                                string
	ListenInterface                         string
	SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION  bool
	SMTP_OPTION_ALLOW_BARE_LINE_FEED_SUBMIT bool
	SMTP_OPTION_MAX_QUEUE_BUFFERED_ITEMS    int
	SMTP_OPTION_MAX_QUEUE_WORKERS           int
	SMTP_OPTION_BACKEND_INTERFACE           string
	Enabled                                 bool
}

type Imap4ServerConfig struct {
	Port                                    uint32
	Hostname                                string
	ListenInterface                         string
	SMTP_OPTION_SINGLE_MESSAGE_PER_SESSION  bool
	SMTP_OPTION_ALLOW_BARE_LINE_FEED_SUBMIT bool
	SMTP_OPTION_MAX_QUEUE_BUFFERED_ITEMS    int
	SMTP_OPTION_MAX_QUEUE_WORKERS           int
	Enabled                                 bool
}

type BackendConfig struct {
	ListenInterface string
	Enabled         bool
}

func NewApplicationConfig() *ApplicationConfig {
	return &ApplicationConfig{}
}

func CreateDefaultAppConfig() *ApplicationConfig {
	c := NewApplicationConfig()
	c.smtpServerConfig = CreateDefaultSmtpServerConfig()
	c.backendConfig = CreateDefaultBackendConfig()
	c.imap4ServerConfig = CreateDefaultImap4ServerConfig()
	return c.PopulateDefault()
}

func CreateDefaultSmtpServerConfig() *SmtpServerConfig {
	c := &SmtpServerConfig{}
	c.Port = 25
	c.ListenInterface = "0.0.0.0"
	c.Hostname = "desktop"
	c.Enabled = true

	// The number of smtp queue workers...
	c.SMTP_OPTION_MAX_QUEUE_WORKERS = 4

	// The number of buffered items allowed in the transport buffer before smtp workers are blocked
	c.SMTP_OPTION_MAX_QUEUE_BUFFERED_ITEMS = 16

	return c
}

func CreateDefaultImap4ServerConfig() *Imap4ServerConfig {
	c := &Imap4ServerConfig{}
	c.Port = 143
	c.ListenInterface = "0.0.0.0"
	c.Hostname = "desktop"
	c.Enabled = true

	// The number of smtp queue workers...
	c.SMTP_OPTION_MAX_QUEUE_WORKERS = 4

	// The number of buffered items allowed in the transport buffer before smtp workers are blocked
	c.SMTP_OPTION_MAX_QUEUE_BUFFERED_ITEMS = 16

	return c
}

func CreateDefaultBackendConfig() *BackendConfig {
	c := &BackendConfig{}

	c.Enabled = true
	c.ListenInterface = "embedded"
	return c
}

func GetAppConfig() *ApplicationConfig {
	return appConfigInstance
}

func (e *ApplicationConfig) PopulateDefault() *ApplicationConfig {
	e.version = "0.0.1 ALPHA"
	return e
}

func (config *ApplicationConfig) GetApplicationVersion() string {
	return config.version
}

func (config *ApplicationConfig) GetLogLevel() uint32 {
	return config.loglevel
}

func (config *ApplicationConfig) GetSmtpServerConfig() *SmtpServerConfig {
	return config.smtpServerConfig
}

func (config *ApplicationConfig) GetImap4ServerConfig() *Imap4ServerConfig {
	return config.imap4ServerConfig
}

func (config *ApplicationConfig) GetBackendConfig() *BackendConfig {
	return config.backendConfig
}
