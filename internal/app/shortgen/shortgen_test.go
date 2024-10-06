package shortgen

import (
	"fmt"
	"testing"
)

// BenchmarkGetShortLink тест генерации ссылки указанного размера
func BenchmarkGetShortLink(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetShortLink(7)
	}
}

func ExampleGetShortLink() {
	short := GetShortLink(7)
	fmt.Println(short)
}
