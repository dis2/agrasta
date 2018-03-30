package agrasta
import "fmt"
import "math/bits"
import "encoding/binary"
import "golang.org/x/crypto/sha3"

type Block [BlockWords]uint64
type Matrix [BlockSize]Block

type State struct {
	rbuf [17]uint64
	rpos int

	sha3.ShakeHash
}

// LU algorithm: We generate only the lower half of LU decomposition.
// Then, we straight proceed to multiply it with the upper - rows of which
// we hallucinate on the fly.
//
// This is still pretty slow to do in Go, but serves as a blueprint for
// reasonably parallelized & unrolled SIMD implementation.
func (s *State) NewMatrix() (res Matrix) {
	// first, generate the lower triangular
	var LT [LTSize]uint64
	var i,j int
	pos := 0
	for i = 0; i < BlockSize; i++ {
		ltlimit := i/64
		for j = 0; j < ltlimit; j++ {
			LT[pos] = s.rand()
			pos++
		}
		one := uint64(1)<<(uint(i)%64)
		LT[pos] = (s.rand() & (one-1)) | one
		pos++
	}

	// now generate the upper left triangular
	// we compute the final LU product on the fly
	for i = BlockSize-1; i >= 0; i-- {
		utlimit := i/64
		// make up our (temporary) upper left row
		var UT [BlockWords]uint64
		for j = 0; j < utlimit; j++ {
			UT[j] = s.rand()
		}
		one := uint64(1)<<(uint(i)%64)
		UT[j] = (s.rand() & (one-1)) | one
		// and multiply it with the other we saved in LT
		pos = 0
		resrow := &res[BlockSize-1-i]
		for j := 0; j < BlockSize; j++ {
			ltlimit := j/64
			var k int
			parity := 0
			for k = 0; k <= ltlimit; k++ {
				parity += bits.OnesCount64(LT[pos] & UT[k])
				pos++
			}
			resrow.OrBit(j, uint64(parity))
		}
	}
	return
}

func (r *Block) GetBit(nn int) uint64 {
	n := uint(nn)
	return (r[n/64] >> (n%64))&1
}

func (r *Block) OrBit(nn int, bit uint64) {
	n := uint(nn)
	r[n/64] |= (bit&1) << (n%64)
}

func (m *Matrix) Dump() {
	for i := 0; i < BlockSize; i++ {
		fmt.Printf("%s\n", m[i].String())
	}
}

// Get matrix rank. Only for testing.
func (m *Matrix) Rank() int {
	mat := *m
	var row int
	for col := 1; col <= BlockSize; col++ {
		if mat[row].GetBit(BlockSize-col) == 0 {
			r := row
			for {
				if !(r < BlockSize && mat[r].GetBit(BlockSize-col) == 0) {
					break
				}
				r++
			}
			if r >= BlockSize {
				continue
			}
			mat[row], mat[r] = mat[r], mat[row]
		}
		for i := row + 1; i < BlockSize; i++ {
			if mat[i].GetBit(BlockSize-col) == 1 {
				for j := 0; j < BlockWords; j++ {
					mat[i][j] ^= mat[row][j]
				}
			}
		}
		row++
		if row == BlockSize {
			break
		}
	}
	return row
}

// Rotate the row rightwards
func (r *Block) Ror(res *Block, n uint) {
	var carry uint64
	for i := 0; i < BlockWords; i++ {
		res[i] = (r[i] >> n) | carry
		carry = r[i] << (BlockSize-n)
	}
	res[0] |= carry
}

// Use nonce and counter to initialize XOF state
func (s *State) SetState(nonce []byte, counter uint64) (c Block) {
	h := sha3.NewShake256()
	binary.Write(h, binary.LittleEndian, counter)
	if nonce != nil {
		h.Write(nonce)
	}
	s.ShakeHash = h
	s.rpos = 0
	return
}

// Generates one PRF block given a key, nonce and counter.
func (s *State) PRF(key *Block) (c Block) {
	// PRF is just the key at first
	c = *key

	r := 0
	var xp1, xp2 Block
	for {
		// Make up round matrix
		m := s.NewMatrix()

		// Multiply vector c with the round matrix
		// TODO perf: coalesce into the LU generator
		var t Block
		for i, mi := range m {
			parity := 0
			for j := 0; j < BlockWords; j++ {
				parity += bits.OnesCount64(mi[j] & c[j])
			}
			t.OrBit(i, uint64(parity))
		}
		c = t

		// add round constants
		for i := 0; i < BlockWords; i++ {
			c[i] ^= s.rand()
		}

		// Substitution, except for the last round
		if r == Rounds {
			break
		}
		r++
		xp1.Ror(&xp1, 1)
		xp2.Ror(&xp2, 2)
		for i := 0; i < BlockWords; i++ {
			c[i] ^= xp2[i] ^ (xp1[i] & xp2[i])
		}
	}

	// Finally, remove the key
	for i := 0; i < BlockWords; i++ {
		c[i] ^= key[i]
	}
	return
}

// Instantiate matrices with nonce/counter and  encrypt one message
func (s *State) Crypt(key, msg *Block, nonce []byte, counter uint64) (c Block) {
	s.SetState(nonce, counter)
	c = s.PRF(key)
	// xor message with the PRF
	for i := 0; i < BlockWords; i++ {
		c[i] ^= msg[i]
	}
	return
}

func main() {
	s := State{}
	s.ShakeHash = sha3.NewShake256()
	mt := s.NewMatrix()
	mt.Dump()
	fmt.Printf("%d\n", mt.Rank())
}


