package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"time"
)

// Create генерация случайного пользователя
func Create(login string, secret string) (int64, error) {
	if login == "" {
		randLogin, err := randString(8)
		if err != nil {
			return 0, errors.New("failed create user")
		}
		login = randLogin
	}

	if secret == "" {
		randSecret, err := randString(15)
		if err != nil {
			return 0, errors.New("failed create user")
		}
		secret = randSecret
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ok, err := repo.CreateUser(ctx, login, secret)
	if err != nil || !ok {
		return 0, errors.New("failed create user")
	}

	id, err := repo.AuthUser(ctx, login, secret)
	if err != nil {
		return 0, errors.New("failed auth new user")
	}

	return id, nil
}

func HasLogin(login string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exists, err := repo.HasUserLogin(ctx, login)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func Authenticate(login string, secret string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id, err := repo.AuthUser(ctx, login, secret)
	if err != nil {
		return 0, errors.New("failed authenticate user")
	}

	return id, nil
}

func randString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.New("error generating random string")
	}

	return hex.EncodeToString(b)[:n], nil
}
