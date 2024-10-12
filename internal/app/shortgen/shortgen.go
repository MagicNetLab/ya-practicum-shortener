package shortgen

import (
	"crypto/rand"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GetShortLink генерация хеша для короткой ссылки
func GetShortLink(n int) string {
	b := make([]rune, n)
	for i := range b {
		val, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[val.Int64()]
	}

	return string(b)
}
