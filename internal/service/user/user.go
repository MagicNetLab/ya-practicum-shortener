package user

import "math/rand"

// Create генерация случайного пользователя
func Create() User {
	return User{
		ID: rand.Intn(999999),
	}
}
