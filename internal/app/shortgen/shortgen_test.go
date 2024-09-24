package shortgen

import "testing"

func BenchmarkGetShortLink(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetShortLink(7)
	}
}
