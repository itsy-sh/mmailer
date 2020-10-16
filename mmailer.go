package mmailer

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Service interface {
	Name() string
	Send(email Email) (res []Response, err error)
	UnmarshalPosthook(body []byte) ([]Posthook, error)
}

type Facade struct {
	Services  []Service
	Selecting SelectStrategy
	Retry     RetryStrategy
}

func New(selecting SelectStrategy, retry RetryStrategy, services ...Service) *Facade {
	return &Facade{
		Services:  services,
		Selecting: selecting,
		Retry:     retry,
	}
}

func (f *Facade) Send(email Email) (res []Response, err error) {
	if len(f.Services) == 0 {
		return nil, errors.New("facade no services to use")
	}

	var service Service
	strategy := f.Selecting
	if strategy == nil {
		strategy = SelectRandom
	}
	service = strategy(f.Services, email)

	if service == nil {
		return nil, errors.New("selected service does not have a mailer associated with it")
	}

	retry := f.Retry
	if retry == nil {
		retry = RetryNone
	}
	return retry(service, email, f.Services)
}

func (f *Facade) UnmarshalPosthook(r *http.Request) (res []Posthook, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	name := strings.ToLower(r.URL.Query().Get("service"))
	for _, s := range f.Services {
		if strings.ToLower(s.Name()) == name {
			return s.UnmarshalPosthook(body)
		}
	}
	return nil, errors.New("could not find a service to unmarshal posthook to")
}
