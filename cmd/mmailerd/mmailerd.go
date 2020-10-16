package main

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/itsy-sh/mmailer"
	"github.com/itsy-sh/mmailer/internal/config"
	"github.com/itsy-sh/mmailer/services/mailjet"
	"github.com/itsy-sh/mmailer/services/mandrill"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var facade *mmailer.Facade

func main() {

	loadServices()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ping"))
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("expect a post request"))
			return
		}

		key := r.URL.Query().Get("key")
		if subtle.ConstantTimeCompare([]byte(key), []byte(config.Get().APIKey)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("not authorized"))
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("could not read body"))
			return
		}

		mail := mmailer.NewEmail()
		err = json.Unmarshal(b, &mail)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("could unmarshal json"))
			return
		}

		res, err := facade.Send(mail)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("could not send email"))
			return
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("could marshal response"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	http.HandleFunc("/posthook", func(w http.ResponseWriter, r *http.Request) {

		key := r.URL.Query().Get("key")
		if subtle.ConstantTimeCompare([]byte(key), []byte(config.Get().PosthookKey)) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("not authorized"))
			return
		}

		resp, err := facade.UnmarshalPosthook(r)
		if err != nil {
			log.Println("[err] ", err)
		}
		if err == nil && len(config.Get().PosthookForward) > 0 {
			go func(hook []mmailer.Posthook) {
				data, err := json.Marshal(hook)
				if err != nil {
					log.Println("[err] could not marshal json, ", err)
					return
				}

				buf := bytes.NewBuffer(data)
				_, err = http.DefaultClient.Post(config.Get().PosthookForward, "application/json", buf)
				if err != nil {
					log.Println("[err] could not post, ", err)
				}
			}(resp)
		}

		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("\n> Send mail by HTTP POST %s/send?key=%s\n\n", config.Get().PublicURL, config.Get().APIKey)

	fmt.Println("Starting server, " + config.Get().HttpInterface)
	log.Fatal(http.ListenAndServe(config.Get().HttpInterface, nil))

}

func loadServices() {
	if len(config.Get().Services) == 0 {
		log.Fatal("Services has to be provide")
	}

	var services []mmailer.Service
	fmt.Println("Services:")
	for _, s := range config.Get().Services {
		parts := strings.Split(s, ":")
		switch strings.ToLower(parts[0]) {
		case "mailjet":
			if len(parts) != 3 {
				log.Println("mailjet api string is not valid,", s)
				continue
			}
			fmt.Printf(" - Mailjet: add the following posthook url  %s/posthook?key=%s&service=mailjet\n", config.Get().PublicURL, config.Get().PosthookKey)
			services = append(services, mailjet.New(parts[1], parts[2]))
		case "mandrill":
			if len(parts) != 2 {
				log.Println("mandrill api string is not valid,", s)
				continue
			}
			fmt.Printf(" - Mandrill: add the following posthook url %s/posthook?key=%s&service=mandrill\n", config.Get().PublicURL, config.Get().PosthookKey)
			services = append(services, mandrill.New(parts[1]))
		case "sendgrid":
			if len(parts) != 2 {
				log.Println("sendgrid api string is not valid,", s)
				continue
			}
			fmt.Printf(" - Sendgrid: add the following posthook url %s/posthook?key=%s&service=sendgrid\n", config.Get().PublicURL, config.Get().PosthookKey)
			services = append(services, mandrill.New(parts[1]))
		}

	}

	if len(services) == 0 {
		log.Fatal("No valid services has to be provide")
	}

	var selects mmailer.SelectStrategy
	fmt.Printf("Select Strategy: ")
	switch config.Get().SelectStrategy {
	case "Weighted":
		fmt.Println("Weighted")
		selects = mmailer.SelectWeighted
	case "RoundRobin":
		fmt.Println("RoundRobin")
		selects = mmailer.SelectRoundRobin()
	case "Random":
		fallthrough
	default:
		fmt.Println("Random")
		selects = mmailer.SelectRandom
	}

	var retry mmailer.RetryStrategy
	fmt.Printf("Retry Strategy:  ")
	switch config.Get().RetryStrategy {
	case "OneOther":
		fmt.Println("OneOther")
		retry = mmailer.RetryOneOther
	case "Each":
		fmt.Println("Each")
		retry = mmailer.RetryEach
	case "Same":
		fmt.Println("Same")
		retry = mmailer.RetrySame
	case "None":
		fallthrough
	default:
		fmt.Println("None")
		retry = mmailer.RetryNone
	}

	facade = mmailer.New(selects, retry, services...)
}
