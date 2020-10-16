package sendgrid

import (
	"errors"
	"github.com/itsy-sh/mmailer"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sendgrid struct {
	apiKey string
}

func (m *Sendgrid) newClient() *sendgrid.Client {
	return sendgrid.NewSendClient(m.apiKey)
}

func New(apiKey string) *Sendgrid {
	return &Sendgrid{
		apiKey: apiKey,
	}
}

func (m *Sendgrid) Name() string {
	return "Sendgrid"
}

func (m *Sendgrid) Send(email mmailer.Email) (res []mmailer.Response, err error) {

	from := mail.NewEmail(email.From.Name, email.From.Email)

	message := mail.NewSingleEmail(from, email.Subject, nil, email.Text, email.Html)

	message.Headers = email.Headers

	for _, a := range email.To {
		message.Personalizations[0].AddTos(&mail.Email{
			Name:    a.Name,
			Address: a.Email,
		})
	}
	for _, a := range email.Cc {
		message.Personalizations[0].AddCCs(&mail.Email{
			Name:    a.Name,
			Address: a.Email,
		})
	}

	response, err := m.newClient().Send(message)
	if err != nil {
		return nil, err
	}

	for _, id := range response.Headers["X-Message-Id"] {
		res = append(res, mmailer.Response{
			Service:   m.Name(),
			MessageId: id,
		})
	}

	return res, nil

}
func (m *Sendgrid) UnmarshalPosthook(body []byte) ([]mmailer.Posthook, error){
	return nil, errors.New("not implemented")
}