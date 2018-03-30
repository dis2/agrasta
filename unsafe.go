// +build !generic

package agrasta

import "unsafe"

func (s *State) rand() uint64 {
	s.rpos--
	if s.rpos < 0 {
		var buf [17*8]byte
		s.ShakeHash.Read(buf[:])
		s.rbuf = *((*[17]uint64)(unsafe.Pointer(&buf)))
		s.rpos = 16
	}
	return s.rbuf[s.rpos]
}


