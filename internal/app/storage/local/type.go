package local

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"os"
)

// StoreEntity структура данных для хранения в кэше
type StoreEntity struct {
	Short     string `json:"short"`
	Link      string `json:"link"`
	UserID    int    `json:"user_id"`
	IsDeleted bool   `json:"is_deleted"`
}

// ErrorLinkNotUnique ошибка в случае попытки сохранения уже существующей в базе ссылки
var ErrorLinkNotUnique = errors.New("link not unique")

// cacheStore кэш для локального хранения данных в файле
type cacheStore struct {
	isInitialized bool
	path          string
}

// Load загрузка данных из файла кэша
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

// Save сохранение данных в файле кэша
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

// IsInitialized проверяет инициализирован кэш или нет
func (cs *cacheStore) IsInitialized() bool {
	return cs.isInitialized
}

// SetPath установка пути до файла кэша
func (cs *cacheStore) SetPath(path string) {
	cs.path = path
}

// SetInitialized установка параметра инициализации кэша
func (cs *cacheStore) SetInitialized(isInitialized bool) {
	cs.isInitialized = isInitialized
}

// linkEntity структура данных для хранения ссылок.
type linkEntity struct {
	userID      int
	shortLink   string
	originalURL string
	isDeleted   bool
}

// store структура с данными ссылок хранимая в памяти приложения.
type store struct {
	isCacheLoaded bool
	store         map[string]linkEntity
}

// PutLink сохранение ссылки в память и кэш
func (s *store) PutLink(ctx context.Context, link string, short string, userID int) error {
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

	cacheStore := getCacheStore()
	_ = cacheStore.Save(l)

	return nil
}

// PutBatchLinksArray пакетное сохранение ссылок в память и кэш
func (s *store) PutBatchLinksArray(ctx context.Context, StoreBatchLinksArray map[string]string, userID int) error {
	rollback := false

	for short, link := range StoreBatchLinksArray {
		if link == "" {
			rollback = true
			break
		}

		err := s.PutLink(ctx, link, short, userID)
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

// GetLink получение ссылки по короткому коду из хранилища в памяти
func (s *store) GetLink(ctx context.Context, short string) (string, bool, error) {
	link, ok := s.store[short]
	if ok {
		return link.originalURL, link.isDeleted, nil
	}

	return "", false, fmt.Errorf("short %s not found", short)
}

// HasShort проверка наличия короткого кода ссылки в хранилище в памяти
func (s *store) HasShort(ctx context.Context, short string) (bool, error) {
	_, ok := s.store[short]

	return ok, nil
}

// GetShort получение ссылки по короткому коду из хранилища в памяти
func (s *store) GetShort(ctx context.Context, link string) (string, error) {
	for k, v := range s.store {
		if v.originalURL == link {
			return k, nil
		}
	}
	return "", fmt.Errorf("short %s not found", link)
}

// GetUserLinks получение всех ссылок пользователя из хранилища в памяти
func (s *store) GetUserLinks(ctx context.Context, userID int) (map[string]string, error) {
	result := make(map[string]string)
	for _, v := range s.store {
		if v.userID == userID {
			result[v.shortLink] = v.originalURL
		}
	}
	return result, nil
}

// DeleteBatchLinksArray пакетное удаление ссылок пользователя из хранилища в памяти и кэше
func (s *store) DeleteBatchLinksArray(ctx context.Context, shorts []string, userID int) error {
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

// Init инициализация хранилища
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

// loadFromFile загрузка данных в память из файла кэша
func (s *store) loadFromFile(filePath string) error {
	cacheStore := getCacheStore()
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
