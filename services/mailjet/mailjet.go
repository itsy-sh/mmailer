package mailjet

import (
	"encoding/json"
	"github.com/itsy-sh/mmailer"
	mj "github.com/mailjet/mailjet-apiv3-go/v3"
	"strings"
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
		TextPart: email.Text,
		HTMLPart: email.Html,
	}

	for k, v := range email.Headers {
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
	if err != nil {
		return nil, err
	}

	for _, rr := range response.ResultsV31 {
		for _, r := range rr.To {
			res = append(res, mmailer.Response{
				Service:   m.Name(),
				MessageId: r.MessageUUID,
				Email:     r.Email,
			})
		}
	}

	return res, nil

}

type Posthook struct {
	Event          string `json:"event"`
	Time           int    `json:"time"`
	MessageID      int64  `json:"MessageID"`
	MessageGUID    string `json:"Message_GUID"`
	Email          string `json:"email"`
	MjCampaignID   int    `json:"mj_campaign_id"`
	MjContactID    int64  `json:"mj_contact_id"`
	Customcampaign string `json:"customcampaign"`
	IP             string `json:"ip"`
	Geo            string `json:"geo"`
	Agent          string `json:"agent"`
	CustomID       string `json:"CustomID"`
	Payload        string `json:"Payload"`
}

func (m *Mailjet) UnmarshalPosthook(body []byte) ([]mmailer.Posthook, error) {
	var hooks []Posthook
	err := json.Unmarshal(body, &hooks)
	if err != nil {
		return nil, err
	}
	var res []mmailer.Posthook
	for _, h := range hooks {
		var event string
		switch strings.ToLower(h.Event) {
		case "sent":
			event = "sent"
		case "open":
			event = "open"
		case "click":
			event = "open"
		case "bounce":
			event = "bounce"
		case "blocked":
			event = "bounce"
		case "spam":
			event = "spam"
		case "unsub":
			event = "spam"
		default:
			event = "unknown"
		}

		res = append(res, mmailer.Posthook{
			Service:   m.Name(),
			MessageId: h.MessageGUID,
			Email:     h.Email,
			Event:     event,
		})
	}
	return res, nil
}
