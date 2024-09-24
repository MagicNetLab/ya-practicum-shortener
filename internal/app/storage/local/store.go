package local

import (
	"errors"
	"fmt"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

type linkEntity struct {
	userID      int
	shortLink   string
	originalURL string
	isDeleted   bool
}

type store struct {
	isCacheLoaded bool
	store         map[string]linkEntity
}

func (s *store) PutLink(link string, short string, userID int) error {
	if link == "" || short == "" {
		args := map[string]interface{}{"link": link, "short": short}
		logger.Error("failed store link: empty link or short", args)
		return errors.New("incorrect params to store link")
	}

	for _, v := range s.store {
		if v.originalURL == link && v.userID == userID {
			return ErrorLinkNotUnique
		}
	}

	l := linkEntity{shortLink: short, originalURL: link, userID: userID, isDeleted: false}
	s.store[short] = l

	cacheStore := GetCacheStore()
	_ = cacheStore.Save(l)

	return nil
}

func (s *store) PutBatchLinksArray(StoreBatchLinksArray map[string]string, userID int) error {
	rollback := false

	for short, link := range StoreBatchLinksArray {
		if link == "" {
			rollback = true
			break
		}

		err := s.PutLink(link, short, userID)
		if err != nil {
			rollback = true
			break
		}
	}

	if rollback {
		for k := range StoreBatchLinksArray {
			delete(s.store, k)
		}

		return errors.New("failed to store batch links array: one of the value is empty")
	}

	return nil
}

func (s *store) GetLink(short string) (string, bool, error) {
	link, ok := s.store[short]
	if ok {
		return link.originalURL, link.isDeleted, nil
	}

	return "", false, fmt.Errorf("short %s not found", short)
}

func (s *store) HasShort(short string) (bool, error) {
	_, ok := s.store[short]

	return ok, nil
}

func (s *store) GetShort(link string) (string, error) {
	for k, v := range s.store {
		if v.originalURL == link {
			return k, nil
		}
	}
	return "", fmt.Errorf("short %s not found", link)
}

func (s *store) GetUserLinks(userID int) (map[string]string, error) {
	result := make(map[string]string)
	for _, v := range s.store {
		if v.userID == userID {
			result[v.shortLink] = v.originalURL
		}
	}
	return result, nil
}

func (s *store) DeleteBatchLinksArray(shorts []string, userID int) error {
	// todo подумать как малой кровью поменять данные в кэш файле
	for _, short := range shorts {
		data, ok := s.store[short]
		if ok && data.userID == userID {
			data.isDeleted = true
			s.store[short] = data
		}
	}

	return nil
}

func (s *store) Init() error {
	conf := config.GetParams()
	if !conf.IsValid() {
		return errors.New("failed init local storage: invalid params")
	}

	if !s.isCacheLoaded {
		err := s.loadFromFile(conf.GetFileStoragePath())
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed load local storage file", args)
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
			s.store[v.Short] = linkEntity{originalURL: v.Link, shortLink: v.Short, userID: v.UserID}
		}
	}

	return nil
}

var Store = store{store: make(map[string]linkEntity, 2)}
