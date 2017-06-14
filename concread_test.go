package concread

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

var implementations []struct {
	name string
	d    impl
}

func init() {
	c, d := newChannel()
	_ = d // Ignoring to signal done for the test/benchmark
	implementations = []struct {
		name string
		d    impl
	}{
		{"mutex", &mutex{}},
		{"channel", c},
		{"atomic", &atom{}},
		{"value", &value{}},
	}
}

func TestSimple(t *testing.T) {
	val := 44

	s := Store{}
	res := s.Read()
	if res != nil {
		t.Error("expected empty content")
	}

	s.Write(val)
	res = s.Read()
	if res != val {
		t.Error("expected same value back")
	}
}

func TestStress(t *testing.T) {
	// Stress test it with multiple combinations of readers and writers, meant to be run under
	// the race detector
	noIt := 100000
	if testing.Short() {
		noIt = 500
	}

	tests := []struct {
		no    int
		ratio int
	}{
		{1, 100},
		{2, 100},
		{10, 100},
		{10, 2},
	}
	for _, i := range implementations {
		for _, te := range tests {
			t.Run(fmt.Sprintf("%v-%v:%v", i.name, te.no, te.ratio), func(t *testing.T) {
				s := i.d
				wg := sync.WaitGroup{}
				for i := 0; i < te.no; i++ {
					wg.Add(1)
					go func() {
						// Randomly decide when we should write and when read
						r := rand.Intn(te.ratio)
						for i := 0; i < noIt; i++ {
							if i%te.ratio == r {
								s.Write(33)
							} else {
								_ = s.Read()
							}
						}
						wg.Done()
					}()
				}
				wg.Wait()
			})
		}
	}
}

// Run a set of benchmarks with different scales for the different
// implementations such that the performance and the scaling can
// be evaluated
func BenchmarkScaling(b *testing.B) {

	type cas struct {
		threads   int
		readratio int
	}
	cases := []cas{
		{1, 1000},
		{2, 1000},
		{4, 1000},
		{1, 10},
		{2, 10},
		{4, 10},
		{1, 1},
		{2, 1},
		{4, 1},
	}

	// Run them implementation by implementation
	for _, i := range implementations {
		for _, c := range cases {
			b.Run(fmt.Sprintf("%v-%v:%v", i.name, c.threads, c.readratio), func(b *testing.B) { benchmark(b, i.d, c.threads, c.readratio) })
		}
	}
}

type impl interface {
	Read() interface{}
	Write(interface{})
}

// What we want to check is how fast read is as the basic measure, it is fine
// if the writes are slower. We do not want to measure the set-up cost but only
// the actual operation cost
func benchmark(b *testing.B, s impl, threads, readratio int) {
	noe := b.N / threads
	wg := sync.WaitGroup{}
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < noe; i++ {
				if i%readratio == 0 {
					s.Write(44)
				} else {
					_ = s.Read()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
