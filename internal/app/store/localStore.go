package store

import "errors"

type local struct {
	store map[string]string
}

func (s *local) PutLink(link string, short string) error {
	s.store[short] = link

	return nil
}

func (s *local) GetLink(short string) (string, error) {
	link, ok := s.store[short]
	if ok {
		return link, nil
	}

	return "", errors.New("short not found")
}

func (s *local) HasShort(short string) bool {
	_, ok := s.store[short]

	return ok
}

var localStore = local{store: make(map[string]string, 2)}
