package mmailer

import (
	"math/rand"
	"sync"
)

type SelectStrategy func([]Service, Email) Service


type WeightService struct {
	Service
	weight uint
}
// Eg. 2 services A with weight 9 and B with Weight 1
// for every 100 messages on average 90 will be sent to A and 10 to B
func NewWeightService(weight uint, service Service) *WeightService {
	return &WeightService{
		Service: service,
		weight:  weight,
	}
}

func SelectRandom(s []Service, _ Email) Service {
	i := rand.Int31n(int32(len(s)))
	return s[i]
}

func SelectRoundRobin() SelectStrategy {
	var i int64
	var mu sync.Mutex
	return func(services []Service, email Email) Service {
		mu.Lock()
		defer mu.Unlock()
		defer func() { i += 1 }()
		l := len(services)
		return services[i%int64(l)]
	}
}



func SelectWeighted(services []Service, e Email) Service {
	var ws []WeightService
	var sum uint
	for _, s := range services {
		w, ok := s.(WeightService)
		if ok {
			sum += w.weight
			ws = append(ws, w)
		}
	}
	if len(ws) == 0 {
		return SelectRandom(services, e)
	}
	rand.Shuffle(len(ws), func(i, j int) {
		ws[i], ws[j] = ws[j], ws[i]
	})

	r := int(rand.Int31n(int32(sum))) + 1

	for _, s := range ws {
		r -= int(s.weight)
		if r <= 0 {
			return s
		}
	}
	return ws[len(ws)-1]
}
