package generic

import (
	"errors"
	"github.com/itsy-sh/mmailer"
	"github.com/itsy-sh/mmailer/internal/smtpx"
	"net/smtp"
	"net/url"
	"strings"
)

type Generic struct {
	smtpUrl *url.URL
}

// smtp://user:pass@smtp.server.com:port
func New(smtpUrl *url.URL) *Generic {
	if smtpUrl.Port() == ""{
		smtpUrl.Host = smtpUrl.Host + ":25"
	}
	return &Generic{
		smtpUrl: smtpUrl,
	}
}
func (m *Generic) Name() string {
	return "Generic smtp " + m.smtpUrl.Host
}

func (m *Generic) Send(email mmailer.Email) (res []mmailer.Response, err error) {


	// Todo think of adding our own message id, in order to be able to track messages..
	message := smtpx.NewMessage()
	for k, v := range email.Headers {
		message.SetHeader(k, v)
	}

	message.SetHeader("From", email.From.String())

	var recp []string
	if len(email.To) > 0 {
		var tos []string
		for _, t := range email.To {
			tos = append(tos, t.String())
			recp = append(recp, t.Email)
		}
		message.SetHeader("To", strings.Join(tos, ", "))
	}
	if len(email.Cc) > 0 {
		var tos []string
		for _, t := range email.To {
			tos = append(tos, t.String())
			recp = append(recp, t.Email)
		}
		message.SetHeader("To", strings.Join(tos, ", "))
	}
	message.SetHeader("Subject", email.Subject)

	if len(email.Text) > 0 {
		message.SetBody("text/plain", email.Text)
	}
	if len(email.Html) > 0 {
		message.SetBody("text/html", email.Html)
	}


	user := m.smtpUrl.User.Username()
	pass, ok := m.smtpUrl.User.Password()
	if !ok {
		return nil, errors.New("could not find password for smtp user")
	}

	msg, err := message.Bytes()
	if err != nil {
		return nil, err
	}
	err = smtp.SendMail(m.smtpUrl.Host, smtp.CRAMMD5Auth(user, pass), email.From.Email, recp, msg)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *Generic) UnmarshalPosthook(body []byte) ([]mmailer.Posthook, error) {
	return nil, errors.New("generic smtp does not have post hooks")
}
