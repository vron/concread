package concread

import (
	"sync"
	"sync/atomic"
)

type value struct {
	mu   sync.Mutex
	data atomic.Value
}

func (s *value) Read() (data interface{}) {
	data = s.data.Load()
	return
}

func (s *value) Write(a interface{}) {
	s.data.Store(a)
}
