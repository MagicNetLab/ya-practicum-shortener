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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed parse cache file", args)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed open local storage file", args)
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
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed serialize cache data", args)
			return err
		}

		_, err = writer.WriteString(string(rowString) + "\n")
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed write local storage file", args)
			return err
		}

		if err := writer.Flush(); err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed flush local storage file", args)
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
