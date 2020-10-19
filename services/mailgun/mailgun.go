package mailgun

import (
	"context"
	"errors"
	"github.com/itsy-sh/mmailer"
	mg "github.com/mailgun/mailgun-go/v4"
)

type Mailgun struct {
	domain        string
	apiKeyPrivate string
}

func (m *Mailgun) newClient() *mg.MailgunImpl {
	return mg.NewMailgun(m.domain, m.apiKeyPrivate)
}

func New(domain, apiKeyPrivate string) *Mailgun {
	return &Mailgun{
		domain:        domain,
		apiKeyPrivate: apiKeyPrivate,
	}
}
func (m *Mailgun) Name() string {
	return "Mailgun"
}

func (m *Mailgun) Send(email mmailer.Email) (res []mmailer.Response, err error) {

	c := m.newClient()
	msg := c.NewMessage(email.From.String(), email.Subject, email.Text)

	for k, v := range email.Headers {
		msg.AddHeader(k, v)
	}

	for _, a := range email.To {
		err := msg.AddRecipient(a.String())
		if err != nil {
			return nil, err
		}
	}
	for _, a := range email.Cc {
		msg.AddCC(a.String())
	}
	msg.SetHtml(email.Html)

	_, id, err := c.Send(context.Background(), msg)
	if err != nil {
		return nil, err
	}
	res = append(res, mmailer.Response{
		Service:   m.Name(),
		MessageId: id,
	})

	return res, nil

}

func (m *Mailgun) UnmarshalPosthook(body []byte) ([]mmailer.Posthook, error) {
	return nil, errors.New("not implemented")
}
