package local

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

type cacheStore struct {
	isInitialized bool
	path          string
}

func (cs *cacheStore) Load() ([]StoreEntity, error) {
	data := make([]StoreEntity, 0)
	file, err := os.OpenFile(cs.path, os.O_RDONLY|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return data, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := StoreEntity{}
		if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
			logger.Log.Errorf("Failed to parse cache file: %s", err)
			return data, err
		}

		data = append(data, row)
	}

	return data, nil
}

func (cs *cacheStore) Save(link linkEntity) error {
	if storeFile.isInitialized {
		file, err := os.OpenFile(storeFile.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger.Log.Errorf("Failed to open local file storage: %s", err)
			return err
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		rowData := StoreEntity{
			Short:     link.shortLink,
			Link:      link.originalURL,
			UserID:    link.userID,
			IsDeleted: link.isDeleted,
		}
		rowString, err := json.Marshal(rowData)
		if err != nil {
			logger.Log.Errorf("Failed to serialize cache data: %s", err)
			return err
		}

		_, err = writer.WriteString(string(rowString) + "\n")
		if err != nil {
			logger.Log.Errorf("Failed to write local file storage: %s", err)
			return err
		}

		if err := writer.Flush(); err != nil {
			logger.Log.Errorf("Failed to flush local file storage: %s", err)
		}
	}

	return nil
}

func (cs *cacheStore) IsInitialized() bool {
	return cs.isInitialized
}

func (cs *cacheStore) SetPath(path string) {
	cs.path = path
}

func (cs *cacheStore) SetInitialized(isInitialized bool) {
	cs.isInitialized = isInitialized
}

var storeFile cacheStore

func GetCacheStore() CacheStoreInterface {
	return &storeFile
}
