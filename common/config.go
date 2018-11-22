package common

var appConfigInstance = CreateDefaultAppConfig()

func init() {
}

type ApplicationConfig struct {
	version          string
	loglevel         uint32
	smtpServerConfig SmtpServerConfig
}

type SmtpServerConfig struct {
	Port            uint32
	Hostname        string
	ListenInterface string
}

func NewApplicationConfig() *ApplicationConfig {
	return &ApplicationConfig{}
}

func CreateDefaultAppConfig() *ApplicationConfig {
	c := NewApplicationConfig()
	c.smtpServerConfig = *CreateDefaultSmtpServerConfig()
	return c.PopulateDefault()
}

func CreateDefaultSmtpServerConfig() *SmtpServerConfig {
	c := &SmtpServerConfig{}
	c.Port = 25
	c.ListenInterface = "0.0.0.0"
	c.Hostname = "desktop"
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
func (config *ApplicationConfig) GetSmtpServerConfig() SmtpServerConfig {
	return config.smtpServerConfig
}
