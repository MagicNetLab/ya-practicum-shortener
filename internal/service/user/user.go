package user

import "math/rand"

type User struct {
	ID int `json:"id"`
}

func Create() User {
	return User{
		ID: rand.Intn(999999),
	}
}
