package agrasta

import (
	"testing"
	"golang.org/x/crypto/sha3"
)

func BenchmarkNewMatrix(b *testing.B) {
	s := State{}
	s.ShakeHash = sha3.NewShake256()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.NewMatrix()
	}
}

func BenchmarkCrypt(b *testing.B) {
	s := State{}
	b.ResetTimer()
	var k Block
	for i := 0; i < b.N; i++ {
		k = s.Crypt(&k, &k, nil, uint64(i))
	}
}


