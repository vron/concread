// Package concread provides a component for efficient multithreaded reads
package concread

// Store supports concurrent reading of the value stored in the store. Note that any
// value returned from Read() may not be modified in any way as this leads to races.
// Read and Write may be called concurrently but the strucutre is intended for cases
// where writes are much more uncommon than reads.
type Store struct {
	atom
}

// Read returns the last value written or nil if not value has been written. The value
// returned must not be modified in any way.
func (s *Store) Read() interface{} {
	return s.atom.Read()
}

// Write stores the particular value in the store such that subsequent reads will return
// the new value until another value is written. Write is safe for concurrent use.
func (s *Store) Write(a interface{}) {
	s.atom.Write(a)
}
