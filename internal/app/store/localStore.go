package store

import (
	"errors"
	"fmt"
	"log"
)

type local struct {
	store map[string]string
}

func (s *local) PutLink(link string, short string) error {
	if link == "" || short == "" {
		log.Printf("Faled store link: empty link(%s) or short(%s)", link, short)
		return errors.New("incorrect params to store link")
	}

	s.store[short] = link

	return nil
}

func (s *local) GetLink(short string) (string, error) {
	link, ok := s.store[short]
	if ok {
		return link, nil
	}

	return "", fmt.Errorf("short %s not found", short)
}

func (s *local) HasShort(short string) bool {
	_, ok := s.store[short]

	return ok
}

var localStore = local{store: make(map[string]string, 2)}
