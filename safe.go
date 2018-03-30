// +build generic

package agrasta

func (s *State) rand() uint64 {
	s.rpos--
	if s.rpos < 0 {
		binary.Read(s.ShakeHash, binary.LittleEndian, &s.rbuf)
		s.rpos = 16
	}
	return s.rbuf[s.rpos]
}


