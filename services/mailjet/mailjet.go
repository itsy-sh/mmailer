package mailjet

import (
	"github.com/itsy-sh/mmailer"
	mj "github.com/mailjet/mailjet-apiv3-go/v3"
)

type Mailjet struct {
	apiKeyPublic  string
	apiKeyPrivate string
}

func (m *Mailjet) newClient() *mj.Client {
	return mj.NewMailjetClient(m.apiKeyPublic, m.apiKeyPrivate)
}

func New(apiKeyPublic, apiKeyPrivate string) *Mailjet {
	return &Mailjet{
		apiKeyPublic:  apiKeyPublic,
		apiKeyPrivate: apiKeyPrivate,
	}
}
func (m *Mailjet) Name() string {
	return "Mailjet"
}

func (m *Mailjet) Send(email mmailer.Email) (res []mmailer.Response, err error) {

	message := mj.InfoMessagesV31{
		Headers: map[string]interface{}{},
		From: &mj.RecipientV31{
			Email: email.From.Email,
			Name:  email.From.Name,
		},
		Subject:  email.Subject,
		TextPart: string(email.Text),
		HTMLPart: string(email.Html),
	}

	for k,v := range email.Headers{
		message.Headers[k] = v
	}

	var to mj.RecipientsV31
	for _, a := range email.To {
		to = append(to, mj.RecipientV31{
			Email: a.Email,
			Name:  a.Name,
		})
	}

	var cc mj.RecipientsV31
	for _, a := range email.Cc {
		cc = append(cc, mj.RecipientV31{
			Email: a.Email,
			Name:  a.Name,
		})
	}
	message.To = &to
	message.Cc = &cc


	messages := mj.MessagesV31{Info: []mj.InfoMessagesV31{message}}
	response, err := m.newClient().SendMailV31(&messages)
	if err != nil{
		return nil, err
	}

	for _, rr := range response.ResultsV31{
		for _, r := range rr.To{
			res = append(res, mmailer.Response{
				Service:   m.Name(),
				MessageId: r.MessageUUID,
			})
		}
	}

	return res, nil

}
