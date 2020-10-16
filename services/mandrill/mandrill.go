package mandrill

import (
	"errors"
	"github.com/itsy-sh/mmailer"
	"github.com/mattbaird/gochimp"
)

type Mandrill struct {
	apiKey string
}

func (m *Mandrill) newClient() (*gochimp.MandrillAPI, error) {
	return gochimp.NewMandrill(m.apiKey)
}

func New(apiKey string) *Mandrill {
	return &Mandrill{
		apiKey: apiKey,
	}
}

func (m *Mandrill) Name() string {
	return "Mandrill"
}

func (m *Mandrill) Send(email mmailer.Email) (res []mmailer.Response, err error){
	c, err := m.newClient()
	if err != nil{
		return nil, err
	}


	var to []gochimp.Recipient
	for _, a := range email.To{
		to = append(to, gochimp.Recipient{
			Email: a.Email,
			Name:  a.Name,
			Type:  "to",
		})
	}

	for _, a := range email.Cc{
		to = append(to, gochimp.Recipient{
			Email: a.Email,
			Name:  a.Name,
			Type:  "cc",
		})
	}

	message := gochimp.Message{
		Headers:   email.Headers,
		FromName:  email.From.Name,
		FromEmail: email.From.Email,
		To:        to,
		Subject:   email.Subject,
		Html:      email.Html,
		Text:      email.Text,
	}

	responses, err := c.MessageSend(message, false)
	if err != nil{
		return nil, err
	}

	for _, r := range responses {
		res = append(res, mmailer.Response{
			Service:   m.Name(),
			MessageId: r.Id,
			Email: r.Email,
		})
	}
	return res, nil

}
func (m *Mandrill) UnmarshalPosthook(body []byte) ([]mmailer.Posthook, error){
	return nil, errors.New("not implemented")
}