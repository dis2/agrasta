// +build agrasta127

package agrasta
import "fmt"

// Number of rounds.
const Rounds = 4

// Block size.
const BlockSize = 127
const BlockWords = 2

// Total size (in 64bit words) of signle matrix triangle.
// The series must go on up until BlockWords
const LTSize = 64 + 64*2

func (r *Block) String() string {
	return fmt.Sprintf("%063b%064b", r[1], r[0])
}


