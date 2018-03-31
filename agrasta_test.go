package agrasta

import (
	"testing"
	"golang.org/x/crypto/sha3"
)

// We have to make sure the rank counting works reliably
func TestMatrixRanker(t *testing.T) {
	s := State{}
	s.ShakeHash = sha3.NewShake256()
	faults := 0

	for i := 0; i < 1000; i++ {
		// Fill the matrix with random data
		var m Matrix
		for i := 0; i < BlockSize; i++ {
			for j := 0; j < BlockWords; j++ {
				m[i][j] = s.rand()
			}
		}
		if m.Rank() != BlockSize {
			faults++
		}
	}
	// It should be around 700, there may be encounter with an outlier case,
	// with fairly low probability - just change hash seed if that happens.
	if faults < 650 || faults > 750 {
		t.Fatalf("Matrix ranking seems off, seen %d/%d faults\n", faults, 1000)
	}
}

func TestNewMatrixRank(t *testing.T) {
	s := State{}
	s.ShakeHash = sha3.NewShake256()
	for i := 0; i < 1000; i++ {
		m := s.NewMatrix()
		if m.Rank() != BlockSize {
			t.Fatal("insufficient rank")
		}
	}
}

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


