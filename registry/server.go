package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	ServerPort  = "3000"
	ServerHost  = "localhost"
	ServicesURL = "http://" + ServerHost + ":" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.registrations = append(r.registrations, reg)
	return nil
}

func (r *registry) remove(reg Registration) error {
	for i := range r.registrations {
		if r.registrations[i].ServiceName == reg.ServiceName {
			r.mutex.Lock()
			defer r.mutex.Unlock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("service %v not found", reg.ServiceName)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Removing service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.remove(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
