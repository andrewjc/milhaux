package smtp

import (
	"fmt"
	"github.com/andrewjc/milhaux/common"
	"github.com/google/uuid"
	"strings"
)

func (s *SmtpCommandProcessor) smtpCommandEstablish(smtpSession *SmtpSession, command commandArgPair) *CommandResponse {
	switch command.commandStr {
	case "HELO": //hello command
		return s.smtpCommandHelo(smtpSession, command)
	case "EHLO": //enhanced hello
		return s.smtpCommandEhlo(smtpSession, command)
	case "STARTTLS": // start authenticated tls session
		return s.smtpCommandStartTls(smtpSession, command)
	default: // wtf is this?
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_BAD_COMMAND_SEQUENCE, "Send EHLO/HELO first!"}
	}
}

func (s *SmtpCommandProcessor) smtpCommandHelo(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.smtpState = SMTP_SERVER_STATE_GOTHELO
	session.stateData[SESSION_DATA_KEY_CLIENT_ID] = commandArg.args

	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("%s Hello", commandArg.args)}
}

func (s *SmtpCommandProcessor) smtpCommandEhlo(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.smtpState = SMTP_SERVER_STATE_GOTHELO
	session.stateData[SESSION_DATA_KEY_CLIENT_ID] = commandArg.args

	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("%s Hello", commandArg.args)}
}

func (s *SmtpCommandProcessor) smtpCommandStartTls(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.smtpState = SMTP_SERVER_STATE_GOTHELO
	session.stateData[SESSION_DATA_KEY_CLIENT_ID] = commandArg.args

	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("%s Hello", commandArg.args)}
}

func (s *SmtpCommandProcessor) smtpCommandDone(smtpSession *SmtpSession, commandLine string) *CommandResponse {
	return &CommandResponse{COMMANDACTION_EXIT, SMTP_COMMAND_STATUS_SERVICE_CLOSING_CHANNEL, "Bye"}
}

func (s *SmtpCommandProcessor) smtpCommandMail(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	switch {
	case strings.HasPrefix(commandArg.args, "FROM:"):
		session.stateData[SESSION_DATA_KEY_MAIL_FROM] = string(commandArg.args[len("FROM:")])
		session.smtpState = SMTP_SERVER_STATE_MAIL

		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, "OK"}
	case strings.HasPrefix(commandArg.args, "TO:"):
		session.stateData[SESSION_DATA_KEY_MAIL_TO] = string(commandArg.args[len("TO:")])
		session.smtpState = SMTP_SERVER_STATE_RCPT
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, "OK"}
	default:
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_COMMAND_NOT_IMPLEMENTED, fmt.Sprintf("Unknown command: %s %s", commandArg.commandStr, commandArg.args)}
	}
}

func (s *SmtpCommandProcessor) smtpCommandReceiveData(session *SmtpSession, commandArg commandArgPair) *CommandResponse {

	session.smtpState = SMTP_SERVER_STATE_DATA
	return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_START_MAIL_INPUT, "Start mail input; end with <CRLF>.<CRLF>"}

}

func (s *SmtpCommandProcessor) smtpCommandBufferData(smtpSession *SmtpSession, commandLine string) *CommandResponse {
	if commandLine == END_DATA_COMMAND_SEQUENCE {
		newMessage := &common.MailMessage{
			To:   smtpSession.stateData[SESSION_DATA_KEY_MAIL_FROM].(string),
			From: smtpSession.stateData[SESSION_DATA_KEY_MAIL_TO].(string),
			Data: smtpSession.receiveDataBuffer.String(),
		}
		s.onSubmitMail(smtpSession, newMessage)
		smtpSession.smtpState = SMTP_SERVER_STATE_DONE
		return &CommandResponse{COMMANDACTION_CONTINUE, SMTP_COMMAND_STATUS_MAIL_ACTION_OK, fmt.Sprintf("OK. Queued for Delivery. Queue id: %s", newMessage.QueueId)}
	} else {
		smtpSession.receiveDataBuffer.WriteString(commandLine)
		return &CommandResponse{COMMANDACTION_NONE, SMTP_COMMAND_STATUS_NONE, ""}
	}
}

func (s *SmtpCommandProcessor) onSubmitMail(session *SmtpSession, message *common.MailMessage) {
	submitQueueMessage := &SmtpServerChannelMessage{
		SMTP_CHANNEL_MESSAGE_QUEUE_SUBMIT,
		message,
	}

	submitQueueMessage.Data.QueueId = uuid.New().String()

	session.smtpMessageChannel <- submitQueueMessage
}
