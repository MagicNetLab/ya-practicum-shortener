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
		logger.Log.Errorf("Failed store link: empty link(%s) or short(%s)", link, short)
		return errors.New("incorrect params to store link")
	}

	s.store[short] = link

	cacheStore := GetCacheStore()
	_ = cacheStore.Save(short, link)

	return nil
}

func (s *store) PutBatchLinksArray(StoreBatchLicksArray map[string]string) error {
	rollback := false

	for k, v := range StoreBatchLicksArray {
		if v == "" {
			rollback = true
			break
		}

		err := s.PutLink(k, v)
		if err != nil {
			rollback = true
			break
		}
	}

	if rollback {
		for k := range StoreBatchLicksArray {
			delete(s.store, k)
		}

		return errors.New("failed to store batch links array: one of the value is empty")
	}

	return nil
}

func (s *store) GetLink(short string) (string, error) {
	link, ok := s.store[short]
	if ok {
		return link, nil
	}

	return "", fmt.Errorf("short %s not found", short)
}

func (s *store) HasShort(short string) (bool, error) {
	_, ok := s.store[short]

	return ok, nil
}

func (s *store) GetShort(link string) (string, error) {
	for k, v := range s.store {
		if v == link {
			return k, nil
		}
	}
	return "", fmt.Errorf("short %s not found", link)
}

func (s *store) Init() error {
	conf := config.GetParams()
	if !conf.IsValid() {
		return errors.New("failed init local storage: invalid params")
	}

	if !s.isCacheLoaded {
		err := s.loadFromFile(conf.GetFileStoragePath())
		if err != nil {
			logger.Log.Errorf("Failed to load local file storage: %s", err)
			return err
		}

		s.isCacheLoaded = true
	}

	return nil
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
