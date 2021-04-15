package lib

import (
	"stock/share/lib/email"
)

type EmailConfig struct {
	Server   string
	Port     string
	Password string
	Addr     string
	Subject  string
	To       string
	Html     string
}

func SendMail(config EmailConfig) error {
	m := email.NewMail()
	m.AddFrom(config.Addr)
	m.AddFromName("haina")
	m.AddTo(config.To)
	m.AddSubject(config.Subject)
	m.AddHTML(config.Html)
	m.AddReplyTo("notification@stock")

	client := email.NewSMTPClient(
		config.Addr,
		config.Password,
		config.Server,
		config.Port)

	err := client.Send(m)
	return err
}
