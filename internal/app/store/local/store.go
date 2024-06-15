package local

import (
	"errors"
	"fmt"
	"log"
)

type store struct {
	store map[string]string
}

func (s *store) PutLink(link string, short string) error {
	if link == "" || short == "" {
		log.Printf("Faled store link: empty link(%s) or short(%s)", link, short)
		return errors.New("incorrect params to store link")
	}

	s.store[short] = link

	return nil
}

func (s *store) GetLink(short string) (string, error) {
	link, ok := s.store[short]
	if ok {
		return link, nil
	}

	return "", fmt.Errorf("short %s not found", short)
}

func (s *store) HasShort(short string) bool {
	_, ok := s.store[short]

	return ok
}

var Store = store{store: make(map[string]string, 2)}
