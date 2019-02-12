package imap

type ImapSessionState uint8

const IMAP_SERVER_BUILD_STRING = "Imap Server"

const (
	IMAP_SESSION_STATE_INVALID = ImapSessionState(iota)
	IMAP_SESSION_STATE_PREAUTH
	IMAP_SESSION_STATE_AUTHOK
	IMAP_SESSION_STATE_CLOSED
)

func (s ImapSessionState) String() string {
	switch s {
	case IMAP_SESSION_STATE_INVALID:
		return "SESSION_STATE_INVALID"
	case IMAP_SESSION_STATE_PREAUTH:
		return "SESSION_STATE_PREAUTH"
	case IMAP_SESSION_STATE_AUTHOK:
		return "SESSION_STATE_AUTHOK"
	case IMAP_SESSION_STATE_CLOSED:
		return "SESSION_STATE_CLOSED"
	}
	return "Unknown"
}

type ImapCommandStatus string

const (
	IMAP_COMMAND_STATUS_NONE                    = ImapCommandStatus("")
	IMAP_COMMAND_STATUS_OK                      = ImapCommandStatus("OK")
	IMAP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED = ImapCommandStatus("ERROR")
)

func (s ImapCommandStatus) String() string {
	switch s {
	case IMAP_COMMAND_STATUS_NONE:
		return ""
	case IMAP_COMMAND_STATUS_OK:
		return "OK"
	case IMAP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED:
		return "NOT_IMPLEMENTED!!"
	}
	return "Unknown"
}

type ImapCommandVerb string

const (
	IMAP_COMMAND_CAPABILITY = ImapCommandVerb("CAPABILITY")
	IMAP_COMMAND_NOOP       = ImapCommandVerb("NOOP")
	IMAP_COMMAND_LOGOUT     = ImapCommandVerb("LOGOUT")
	IMAP_COMMAND_LOGIN      = ImapCommandVerb("LOGIN")
)

func (s ImapCommandVerb) String() string {
	switch s {
	case IMAP_COMMAND_CAPABILITY:
		return "CAPABILITY"
	case IMAP_COMMAND_NOOP:
		return "NOOP"
	case IMAP_COMMAND_LOGOUT:
		return "LOGOUT"
	}
	return "Unknown"
}
