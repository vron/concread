package concread

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type atom struct {
	mu   sync.RWMutex
	data unsafe.Pointer
}

func (s *atom) Read() (data interface{}) {
	p := atomic.LoadPointer(&s.data)
	d := ((*interface{})(p))
	if d != nil {
		data = *d
	}
	return data
}

func (s *atom) Write(a interface{}) {
	if a == nil {
		panic("cannot store nil")
	}
	atomic.StorePointer(&s.data, unsafe.Pointer(&a))
}
