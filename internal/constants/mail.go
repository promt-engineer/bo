package constants

import "text/template"

const (
	MailNotifyUserSubject     = "Backoffice account"
	MailNotifyUserTemplateRaw = `Welcome to: {{.FrontURL}}
Your login: {{.Login}}
Your password: {{.Password}}`
	MailResetUserPasswordRaw = `Password reset page: {{.ResetPasswordURL}}
Your token: {{.Token}}`
)

var MailNotifyUserTemplate *template.Template
var MailResetPasswordTemplate *template.Template

type MailNotifyUserContent struct {
	FrontURL, Login, Password string
}

type MailResetPasswordContent struct {
	ResetPasswordURL, Token string
}

func init() {
	var err error
	if MailNotifyUserTemplate, err = template.New("simulation-txt").Parse(MailNotifyUserTemplateRaw); err != nil {
		panic(err)
	}
	if MailResetPasswordTemplate, err = template.New("reset-txt").Parse(MailResetUserPasswordRaw); err != nil {
		panic(err)
	}
}
