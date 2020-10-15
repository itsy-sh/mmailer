package mmailer

import (
	"errors"
	"fmt"
)

type RetryStrategy func(serviceToUse Service, email Email, backupServices []Service) (res []Response, err error)

func RetryNone(s Service, e Email, _ []Service) (res []Response, err error){
	return s.Send(e)
}


func RetryEach(s Service, e Email, services []Service) (res []Response, err error){
	res, err = s.Send(e)
	if err == nil{
		return res, nil
	}

	var acc string = err.Error()
	for _, ss := range services{
		res, err = ss.Send(e)
		if err == nil{
			return res, nil
		}
		acc = fmt.Sprintf("%s: %s", err.Error(), acc)
	}
	return nil, errors.New(acc)
}

func RetryOneOther(s Service, e Email, services []Service) (res []Response, err error){
	res, err = s.Send(e)
	if err == nil{
		return res, nil
	}
	for _, ss := range services{
		if s.Name() == ss.Name(){
			continue
		}
		return ss.Send(e)
	}
	return nil, err
}

func RetrySame(s Service, e Email, services []Service) (res []Response, err error){
	res, err = s.Send(e)
	if err == nil{
		return res, nil
	}
	return s.Send(e)
}