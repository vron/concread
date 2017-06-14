package concread

func newChannel() (c *channel, close func()) {
	c = &channel{
		read:  make(chan interface{}, 32),
		write: make(chan interface{}, 1),
	}
	done := make(chan struct{})
	go func() {
		for {
			select {
			case c.read <- c.data:
				continue
			case c.data = <-c.write:
				continue
			case <-done:
				return
			}
		}
	}()
	return c, func() {
		done <- struct{}{}
	}
}

type channel struct {
	read  chan interface{}
	write chan interface{}
	data  interface{}
}

func (s *channel) Read() interface{} {

	return <-s.read
}

func (s *channel) Write(a interface{}) {
	s.write <- a
}
