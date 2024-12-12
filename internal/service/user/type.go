package user

// User объект пользователя
type User struct {
	ID     int64  `json:"id"`
	Login  string `json:"login"`
	Secret string `json:"secret"`
}
