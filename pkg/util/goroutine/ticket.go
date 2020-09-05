package goroutine

import "sync"

type ConcurrentControl struct {
	concurrentNum       int
	closeOnce           sync.Once
	concurrentContainer chan struct{}
}

func NewConcurrentControl(concurrentNum int) *ConcurrentControl {
	if concurrentNum <= 0 {
		panic("concurrentNum should more than 1 or equal 1")
	}
	return &ConcurrentControl{
		concurrentNum:       concurrentNum,
		concurrentContainer: make(chan struct{}, concurrentNum),
	}
}

func (c *ConcurrentControl) GetGoroutine() {
	c.concurrentContainer <- struct{}{}
}

func (c *ConcurrentControl) ReleaseConnGoroutine() {
	<-c.concurrentContainer
}

func (c *ConcurrentControl) CloseGoroutine() {
	c.closeOnce.Do(func() {
		close(c.concurrentContainer)
	})
}
