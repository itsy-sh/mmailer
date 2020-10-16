package mmailer

type Address struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Email struct {
	Headers map[string]string `json:"headers"`
	From    Address           `json:"from"`
	To      []Address         `json:"to"`
	Cc      []Address         `json:"cc"`
	Subject string            `json:"subject"`
	Text    string            `json:"text"`
	Html    string            `json:"html"`
}

func NewEmail() Email {
	return Email{
		Headers: map[string]string{},
	}
}

type Response struct {
	Service   string `json:"service"`
	MessageId string `json:"message_id"`
	Email     string `json:"email"`
}

type Posthook struct {
	Service   string `json:"service"`
	MessageId string `json:"message_id"`
	Email     string `json:"email"`
	Event     string `json:"event"`
}
