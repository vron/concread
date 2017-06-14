package concread

import "sync"

type mutex struct {
	mu   sync.RWMutex
	data interface{}
}

func (s *mutex) Read() interface{} {
	s.mu.RLock()
	d := s.data
	s.mu.RUnlock()
	return d
}

func (s *mutex) Write(a interface{}) {
	s.mu.Lock()
	s.data = a
	s.mu.Unlock()
}
