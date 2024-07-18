package cache

import (
	"time"
)

var (
	checkValidInterval = time.Minute * 10
)

type Pool struct {
	data map[string]*Cache
	size int
}

func NewPool(size int) *Pool {
	c := &Pool{
		data: make(map[string]*Cache),
		size: size,
	}
	go c.checkValid()
	return c
}

func (p *Pool) Key(key string) *Cache {
	if _, ok := p.data[key]; !ok {
		p.data[key] = newCache(p.size)
	}
	return p.data[key]
}

func (p *Pool) checkValid() {
	ticker := time.NewTicker(checkValidInterval)
	for range ticker.C {
		for key, value := range p.data {
			if !value.isValid() {
				delete(p.data, key)
			}
		}
	}
	defer ticker.Stop()
}
