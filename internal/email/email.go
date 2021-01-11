package email

import (
	"net/smtp"
	"strconv"

	"github.com/jordan-wright/email"
	"github.com/kornelkabele/watchdog/internal/cfg"
)

func SendEmail(subj string, body string, attachments []string) error {
	e := email.NewEmail()
	e.From = cfg.SMTP.Sender
	e.To = []string{cfg.SMTP.Receiver}
	e.Subject = subj
	e.Text = []byte(body)
	for _, v := range attachments {
		e.AttachFile(v)
	}
	if err := e.Send(cfg.SMTP.Host+":"+strconv.Itoa(cfg.SMTP.Port), smtp.PlainAuth("", cfg.SMTP.User, cfg.SMTP.Pass, cfg.SMTP.Host)); err != nil {
		return err
	}
	return nil
}
