package imap

import (
	"github.com/andrewjc/milhaux/common"
	log "github.com/sirupsen/logrus"
)

type Imap4Server interface {
}

type Imap4Server_Impl struct {
	config *common.ApplicationConfig
}

func NewIMap4Server(config *common.ApplicationConfig) Imap4Server {
	log.Debug("Creating new imap4 server instance...")

	imapSvr := Imap4Server_Impl{config}
	return imapSvr
}
