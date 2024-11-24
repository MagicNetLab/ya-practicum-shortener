package memory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/exp/slices"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

// GetStore возвращает ссылку на объект хранилища
func GetStore() *Store {
	return &Store{}
}

// Store объект хранилища данных
type Store struct {
	data map[string]linkEntity
	file string
}

// PutLink сохранение ссылки пользователя в хранилище.
func (s *Store) PutLink(ctx context.Context, link string, short string, userID int) error {
	if link == "" || short == "" {
		args := map[string]interface{}{"link": link, "short": short}
		logger.Error("failed store link: empty link or short", args)
		return errors.New("incorrect params to store link")
	}

	for _, v := range s.data {
		if v.originalURL == link && v.userID == userID {
			return ErrorLinkNotUnique
		}
	}

	l := linkEntity{shortLink: short, originalURL: link, userID: userID, isDeleted: false}
	s.data[short] = l

	return nil
}

// PutBatchLinksArray сохранение пакета ссылок пользователя в хранилище.
func (s *Store) PutBatchLinksArray(ctx context.Context, StoreBatchLinksArray map[string]string, userID int) error {
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
			delete(s.data, k)
		}

		return errors.New("failed to store batch links array: one of the value is empty")
	}

	return nil
}

// GetLink получение оригинальной ссылки по короткому хэшу.
func (s *Store) GetLink(ctx context.Context, short string) (string, bool, error) {
	link, ok := s.data[short]
	if ok {
		return link.originalURL, link.isDeleted, nil
	}

	return "", false, fmt.Errorf("short %s not found", short)
}

// HasShort проверка наличия коротко ссылки в хранилище
func (s *Store) HasShort(ctx context.Context, short string) (bool, error) {
	_, ok := s.data[short]

	return ok, nil
}

// GetShort получение короткой ссылки из хранилища для оригинальной ссылки
func (s *Store) GetShort(ctx context.Context, link string) (string, error) {
	for k, v := range s.data {
		if v.originalURL == link {
			return k, nil
		}
	}
	return "", fmt.Errorf("short %s not found", link)
}

// GetUserLinks получение всех ссылок пользователя из хранилища
func (s *Store) GetUserLinks(ctx context.Context, userID int) (map[string]string, error) {
	result := make(map[string]string)
	for _, v := range s.data {
		if v.userID == userID {
			result[v.shortLink] = v.originalURL
		}
	}
	return result, nil
}

// DeleteBatchLinksArray пометка массива ссылок пользователя как удаленных
func (s *Store) DeleteBatchLinksArray(ctx context.Context, shorts []string, userID int) error {
	for _, short := range shorts {
		data, ok := s.data[short]
		if ok && data.userID == userID {
			data.isDeleted = true
			s.data[short] = data
		}
	}

	return nil
}

// GetLinksCount возвращает количество сокращенных ссылок в системе
func (s *Store) GetLinksCount(ctx context.Context) (int, error) {
	var count int
	for _, v := range s.data {
		if !v.isDeleted {
			count++
		}
	}

	return count, nil
}

// GetUsersCount возвращает количество пользователей в системе
func (s *Store) GetUsersCount(ctx context.Context) (int, error) {
	var users []int
	for _, v := range s.data {
		if !slices.Contains(users, v.userID) {
			users = append(users, v.userID)
		}
	}

	return len(users), nil
}

// Initialize инициализация хранилища
func (s *Store) Initialize(config *config.Configurator) error {
	s.data = make(map[string]linkEntity, 10)
	if s.file = config.GetFileStoragePath(); s.file != "" {
		err := s.loadFromFile()
		if err != nil {
			args := map[string]interface{}{"error": err.Error(), "fileName": s.file}
			logger.Error("error loading data from file. cache not initialize", args)
			s.file = ""
		}
	}

	return nil
}

// Close Закрывает хранилище
func (s *Store) Close() error {
	err := s.syncCache()
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("filed sync data to cache", args)
		return err
	}

	return nil
}

// syncCache запись данных в кэш
func (s *Store) syncCache() error {
	if s.file != "" {
		file, err := os.OpenFile(s.file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			args := map[string]interface{}{"error": err.Error()}
			logger.Error("failed open local storage file", args)
			return err
		}
		defer file.Close()

		writer := bufio.NewWriter(file)

		for _, link := range s.data {
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
				return err
			}
		}
	}

	return nil
}

// loadFromFile загрузка дынных из кэша на диске
func (s *Store) loadFromFile() error {
	if s.file != "" {
		data := make([]StoreEntity, 0)
		f, err := os.OpenFile(s.file, os.O_RDONLY|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			row := StoreEntity{}
			if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
				args := map[string]interface{}{"error": err.Error()}
				logger.Error("failed parse cache file", args)
				return err
			}

			data = append(data, row)
		}

		for _, v := range data {
			s.data[v.Short] = linkEntity{originalURL: v.Link, shortLink: v.Short, userID: v.UserID}
		}

		return nil
	}

	return nil
}
