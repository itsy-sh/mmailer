package mmailer

type Address struct {
	Name string
	Email string
}

type Email struct {
	Headers map[string]string
	Subject string
	From Address
	To []Address
	Cc []Address
	//Bcc []Address  // Todo has to be validated
	Text []byte
	Html []byte
}

func NewEmail() Email{
	return Email{
		Headers: map[string]string{},
	}
}



type Response struct {
	Service string
	MessageId string
}