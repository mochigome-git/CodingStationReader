// pkg/util/counter.go
package util

import "sync"

type Counter struct {
	sync.Mutex
	count int
}

func (c *Counter) Add(n int) {
	c.Lock()
	defer c.Unlock()
	// If the insert is successful, increment the count
	c.count += n
}

func (c *Counter) Count() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}

func (c *Counter) Reset() {
	c.Lock()
	defer c.Unlock()
	c.count = 0
}
