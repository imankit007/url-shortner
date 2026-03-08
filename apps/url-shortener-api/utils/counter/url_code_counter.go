package counter

import "sync/atomic"

type URLCodeCounter struct {
	value atomic.Uint64
}

func NewURLCodeCounter() *URLCodeCounter {
	urlCodeCounter := &URLCodeCounter{}
	urlCodeCounter.value.Store(123456)
	return urlCodeCounter
}

func (c *URLCodeCounter) Next() uint64 {
	return c.value.Add(1) - 1
}
