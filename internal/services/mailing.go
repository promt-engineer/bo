package services

import (
	"backoffice/internal/constants"
	"backoffice/pkg/mailgun"
	"bytes"
)

type MailingService struct {
	mailgun          *mailgun.Client
	frontURL         string
	sendEmail        string
	resetPasswordURL string
}

func NewMailingService(mailgunClient *mailgun.Client, frontURL, sendEmail, resetPasswordURL string) *MailingService {
	return &MailingService{mailgun: mailgunClient, frontURL: frontURL, sendEmail: sendEmail, resetPasswordURL: resetPasswordURL}
}

func (s *MailingService) NotifyUserEmail(email, login, password string) error {
	buf := bytes.NewBufferString("")
	err := constants.MailNotifyUserTemplate.
		Execute(buf, constants.MailNotifyUserContent{FrontURL: s.frontURL, Login: login, Password: password})

	if err != nil {
		return err
	}

	s.mailgun.Send(constants.MailNotifyUserSubject, email, s.sendEmail, buf.String(), nil, nil)

	return nil
}

func (s *MailingService) ResetUserPassword(email, token string) error {
	buf := bytes.NewBufferString("")
	err := constants.MailResetPasswordTemplate.
		Execute(buf, constants.MailResetPasswordContent{ResetPasswordURL: s.resetPasswordURL, Token: token})
	if err != nil {
		return err
	}

	s.mailgun.Send(constants.MailNotifyUserSubject, email, s.sendEmail, buf.String(), nil, nil)

	return nil
}
