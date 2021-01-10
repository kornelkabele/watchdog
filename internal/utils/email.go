package utils

import (
	"net/smtp"
	"strconv"

	"github.com/jordan-wright/email"
)

func sendEmail(subj string, body string, attachments []string, cfg ConfigSMTP) error {
	e := email.NewEmail()
	e.From = cfg.Sender
	e.To = []string{cfg.Receiver}
	e.Subject = subj
	e.Text = []byte(body)
	for _, v := range attachments {
		e.AttachFile(v)
	}
	if err := e.Send(cfg.Host+":"+strconv.Itoa(cfg.Port), smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)); err != nil {
		return err
	}
	return nil
}
