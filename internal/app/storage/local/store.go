package local

import (
	"errors"
	"fmt"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

// todo use sync.RWMutex???
type store struct {
	isCacheLoaded bool
	store         map[string]string
}

func (s *store) PutLink(link string, short string) error {
	if link == "" || short == "" {
		logger.Log.Errorf("Failed store link: empty link(%s) or short(%s)", link, short)
		return errors.New("incorrect params to store link")
	}

	s.store[short] = link

	cacheStore := GetCacheStore()
	_ = cacheStore.Save(short, link)

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
		err := s.loadFromFile(conf.GetFileStoragePath())
		if err != nil {
			logger.Log.Errorf("Failed to load local file storage: %s", err)
			return
		}

		s.isCacheLoaded = true
	}
}

func (s *store) loadFromFile(filePath string) error {
	cacheStore := GetCacheStore()
	if filePath != "" {
		cacheStore.SetPath(filePath)
		cacheStore.SetInitialized(true)
		data, err := cacheStore.Load()
		if err != nil {
			return err
		}

		for _, v := range data {
			s.store[v.Short] = v.Link
		}
	}

	return nil
}

var Store = store{store: make(map[string]string, 2)}
