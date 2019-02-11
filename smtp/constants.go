package smtp

type SmtpSessionState uint8

const (
	SMTP_SERVER_BUILD_STRING = "Milhaux SMTP Server"
)

const (
	SMTP_SESSION_STATE_INVALID = SmtpSessionState(iota)
	SMTP_SESSION_STATE_PREAUTH
	SMTP_SESSION_STATE_AUTHOK
	SMTP_SESSION_STATE_MAIL
	SMTP_SESSION_STATE_RCPT
	SMTP_SESSION_STATE_DATA
	SMTP_SESSION_STATE_SUBMIT
	SMTP_SESSION_STATE_CLOSED
)

func (s SmtpSessionState) String() string {
	switch s {
	case SMTP_SESSION_STATE_INVALID:
		return "SESSION_STATE_INVALID"
	case SMTP_SESSION_STATE_PREAUTH:
		return "SESSION_STATE_PREAUTH"
	case SMTP_SESSION_STATE_AUTHOK:
		return "SESSION_STATE_AUTHOK"
	case SMTP_SESSION_STATE_MAIL:
		return "SESSION_STATE_MAIL"
	case SMTP_SESSION_STATE_RCPT:
		return "SESSION_STATE_RCPT"
	case SMTP_SESSION_STATE_DATA:
		return "SESSION_STATE_DATA"
	case SMTP_SESSION_STATE_SUBMIT:
		return "SESSION_STATE_SUBMIT"
	case SMTP_SESSION_STATE_CLOSED:
		return "SESSION_STATE_CLOSED"
	}
	return "Unknown"
}

type SmtpCommandStatus uint16

const (
	SMTP_COMMAND_STATUS_NONE                     = SmtpCommandStatus(0)
	SMTP_COMMAND_STATUS_SERVICE_READY            = SmtpCommandStatus(220)
	SMTP_COMMAND_STATUS_MAIL_ACTION_OK           = SmtpCommandStatus(250)
	SMTP_COMMAND_STATUS_SERVICE_CLOSING_CHANNEL  = SmtpCommandStatus(221)
	SMTP_COMMAND_STATUS_START_MAIL_INPUT         = SmtpCommandStatus(354)
	SMTP_COMMAND_STATUS_REQUESTED_ACTION_ABORTED = SmtpCommandStatus(451)
	SMTP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED  = SmtpCommandStatus(502)
	SMTP_COMMAND_STATUS_BAD_COMMAND_SEQUENCE     = SmtpCommandStatus(503)
)

func (s SmtpCommandStatus) String() string {
	switch s {
	case SMTP_COMMAND_STATUS_NONE:
		return "STATUS_NONE"
	case SMTP_COMMAND_STATUS_SERVICE_READY:
		return "STATUS_SERVICE_READY"
	case SMTP_COMMAND_STATUS_MAIL_ACTION_OK:
		return "STATUS_MAIL_ACTION_OK"
	case SMTP_COMMAND_STATUS_SERVICE_CLOSING_CHANNEL:
		return "STATUS_SERVICE_CLOSING_CHANNEL"
	case SMTP_COMMAND_STATUS_START_MAIL_INPUT:
		return "STATUS_START_MAIL_INPUT"
	case SMTP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED:
		return "STATUS_COMMAND_NOT_IMPLEMENTED"
	case SMTP_COMMAND_STATUS_BAD_COMMAND_SEQUENCE:
		return "STATUS_BAD_COMMAND_SEQUENCE"
	}
	return "Unknown"
}

type SmtpCommandVerb string

const (
	SMTP_COMMAND_MAIL = SmtpCommandVerb("CAPABILITY")

	SMTP_COMMAND_RCPT = SmtpCommandVerb("RCPT")
	SMTP_COMMAND_DATA = SmtpCommandVerb("DATA")
	SMTP_COMMAND_DONE = SmtpCommandVerb("QUIT")
)

func (s SmtpCommandVerb) String() string {
	switch s {
	case SMTP_COMMAND_MAIL:
		return "MAIL"
	case SMTP_COMMAND_RCPT:
		return "RCPT"
	case SMTP_COMMAND_DATA:
		return "DATA"
	case SMTP_COMMAND_DONE:
		return "DONE"
	}
	return "Unknown"
}
