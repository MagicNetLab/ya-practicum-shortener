package local

import (
	"errors"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

type store struct {
	isCacheLoaded bool
	store         map[string]string
}

func (s *store) PutLink(link string, short string) error {
	if link == "" || short == "" {
		logger.Log.Errorf("Faled store link: empty link(%s) or short(%s)", link, short)
		return errors.New("incorrect params to store link")
	}

	s.store[short] = link

	_ = storeFile.Save(short, link)

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

func (s *store) Init() {
	conf := config.GetParams()
	if conf.IsValid() && !s.isCacheLoaded {
		fileStoragePath := conf.GetFileStoragePath()
		if fileStoragePath != "" {
			storeFile.SetPath(fileStoragePath)
			storeFile.SetInitialized(true)
			data, err := storeFile.Load()
			if err != nil {
				logger.Log.Errorf("Failed to load local file storage: %s", err)
				return
			}

			for _, v := range data {
				s.store[v.Short] = v.Link
			}
		}
		s.isCacheLoaded = true
	}
}

var Store = store{store: make(map[string]string, 2)}
