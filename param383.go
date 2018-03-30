// +build agrasta383

package agrasta
import "fmt"

// Number of rounds.
const Rounds = 6

// Block size.
const BlockSize = 383
const BlockWords = 6

// Total size (in 64bit words) of signle matrix triangle.
// The series must go on up until BlockWords
const LTSize = 64 + 64*2 + 64*3 + 64*4 + 64*5 + 64*6

func (r *Block) String() string {
	return fmt.Sprintf("%063b%064b%064b%064b%064b%064b", r[5], r[4], r[3], r[2], r[1], r[0])
}


