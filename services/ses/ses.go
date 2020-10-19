package ses

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsese "github.com/aws/aws-sdk-go/service/ses"
	"github.com/itsy-sh/mmailer"
	"net/url"
)

type SES struct {
	smtpUrl *url.URL
}

// smtp://user:pass@smtp.server.com:port
func New(smtpUrl *url.URL) *SES {
	if smtpUrl.Port() == "" {
		smtpUrl.Host = smtpUrl.Host + ":25"
	}
	return &SES{
		smtpUrl: smtpUrl,
	}
}
func (m *SES) Name() string {
	return "SES"
}

func (m *SES) Send(email mmailer.Email) (res []mmailer.Response, err error) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create an SES session.
	svc := awsese.New(sess)

	input := &awsese.SendEmailInput{
		Destination: &awsese.Destination{
			ToAddresses: []*string{},
			CcAddresses: []*string{},
		},
		Message: &awsese.Message{
			Body: &awsese.Body{
				Html: &awsese.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(email.Html),
				},
				Text: &awsese.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(email.Text),
				},
			},
			Subject: &awsese.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(email.From.String()),
	}

	for _, t := range email.To {
		input.Destination.ToAddresses = append(input.Destination.ToAddresses, aws.String(t.String()))
	}
	for _, t := range email.Cc {
		input.Destination.CcAddresses = append(input.Destination.CcAddresses, aws.String(t.String()))
	}

	_, err = svc.SendEmail(input)

	return nil, err
}

func (m *SES) UnmarshalPosthook(body []byte) ([]mmailer.Posthook, error) {
	return nil, errors.New("generic smtp does not have post hooks")
}
