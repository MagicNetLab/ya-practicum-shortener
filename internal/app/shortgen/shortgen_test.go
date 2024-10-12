package shortgen

import (
	"testing"
)

// BenchmarkGetShortLink тест генерации ссылки указанного размера
func BenchmarkGetShortLink(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetShortLink(7)
	}
}
