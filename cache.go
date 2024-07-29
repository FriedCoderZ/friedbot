package friedbot

import (
	"sync"
	"time"
)

var (
	expDuration   = time.Minute * 15
	checkInterval = time.Minute * 1
)

type Ring struct {
	*sync.RWMutex
	data    []Event
	expTime time.Time
	size    int
	index   int
	length  int
}

func newRing(size int) *Ring {
	return &Ring{
		RWMutex: &sync.RWMutex{},
		data:    []Event{},
		expTime: time.Now().Add(expDuration),
		size:    size,
		index:   -1,
		length:  0,
	}
}

// Ensure all public methods acquire a lock before executing.
func (r *Ring) isValid() bool {
	r.RLock()
	defer r.RUnlock()
	if len(r.data) == 0 {
		return false
	}
	currentTime := time.Now()
	if currentTime.After(r.expTime) {
		return false
	}
	return true
}

func (r *Ring) Len() int {
	r.RLock()
	defer r.RUnlock()
	if r.isValid() {
		return r.length
	} else {
		return 0
	}
}

func (r *Ring) refreshExp(expDuration time.Duration) {
	currentTime := time.Now()
	r.expTime = currentTime.Add(expDuration)
}

func (r *Ring) Add(event Event) {
	r.Lock()
	defer r.Unlock()
	r.index = (r.index + 1) % r.size
	if len(r.data) < r.size {
		r.data = append(r.data, event)
		r.length++
	} else {
		r.data[r.index] = event
	}
}

func (r *Ring) All() (values []Event) {
	if r == nil {
		return nil
	}
	r.RLock()
	defer r.RUnlock()
	r.refreshExp(expDuration)
	return append(r.data[r.index:len(r.data)], r.data[:r.index]...)
}

func (r *Ring) Top() Event {
	if r == nil {
		return nil
	}
	r.RLock()
	defer r.RUnlock()
	r.refreshExp(expDuration)
	return r.data[r.index]
}

func (r *Ring) Limit(length int) (values []Event) {
	if r == nil {
		return nil
	}
	r.RLock()
	defer r.RUnlock()
	if length < len(r.data) {
		return r.data
	}
	var result []Event
	if r.index > length {
		result = r.data[r.index-length : r.index]
	} else {
		result = append(r.data[len(r.data)-length+r.index:len(r.data)], r.data[:r.index]...)
	}
	r.refreshExp(expDuration)
	return result
}

type Pool struct {
	data map[string]*Ring
	size int
}

func NewPool(size int) *Pool {
	c := &Pool{
		data: make(map[string]*Ring),
		size: size,
	}
	go c.checkValid()
	return c
}

func (p *Pool) Key(key string) *Ring {
	if _, ok := p.data[key]; !ok {
		p.data[key] = newRing(p.size)
	}
	return p.data[key]
}

func (p *Pool) checkValid() {
	ticker := time.NewTicker(checkInterval)
	for range ticker.C {
		for key, value := range p.data {
			if !value.isValid() {
				delete(p.data, key)
			}
		}
	}
	defer ticker.Stop()
}

type Cache struct {
	GroupEventsPool   *Pool
	PrivateEventsPool *Pool
	NoticeEventsPool  *Pool
	MetaEventsPool    *Pool
	RequestEventsPool *Pool
}

var service *Cache

func NewService() *Cache {
	service = &Cache{}
	return service
}

func GetService() *Cache {
	return service
}

const (
	CacheSizeLittle = 8
	CacheSizeSmall  = 32
	CacheSizeNormal = 128
	CacheSizeLarge  = 512
	CacheSizeHuge   = 2048
)

func (s *Cache) Start() {
	s.GroupEventsPool = NewPool(CacheSizeLarge)
	s.PrivateEventsPool = NewPool(CacheSizeNormal)
	s.NoticeEventsPool = NewPool(CacheSizeSmall)
	s.MetaEventsPool = NewPool(CacheSizeLittle)
	s.RequestEventsPool = NewPool(CacheSizeLittle)
}
