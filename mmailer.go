package mmailer

import (
	"errors"
)

type Service interface {
	Name() string
	Send(email Email) (res []Response, err error)
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
	if retry == nil{
		retry = RetryNone
	}
	return retry(service, email, f.Services)
}
