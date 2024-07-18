package cache

import (
	"sync"
	"time"

	"github.com/FriedCoderZ/friedbot/models"
)

var (
	expDuration = time.Minute * 15
)

type Cache struct {
	sync.RWMutex
	data    []models.Event
	expTime time.Time
	index   int
	length  int
}

func newCache(size int) *Cache {
	return &Cache{
		RWMutex: sync.RWMutex{},
		data:    make([]models.Event, size),
		expTime: time.Now().Add(expDuration),
		index:   0,
		length:  0,
	}
}

// Ensure all public methods acquire a lock before executing.
func (q *Cache) isValid() bool {
	q.RLock()
	defer q.RUnlock()
	if len(q.data) == 0 {
		return false
	}
	currentTime := time.Now()
	if currentTime.After(q.expTime) {
		return false
	}
	return true
}

func (q *Cache) Len() int {
	q.RLock()
	defer q.RUnlock()
	if q.isValid() {
		return q.length
	} else {
		return 0
	}
}

func (q *Cache) refreshExp(expDuration time.Duration) {
	q.Lock()
	defer q.Unlock()
	currentTime := time.Now()
	q.expTime = currentTime.Add(expDuration)
}

func (q *Cache) Add(data *models.Event) {
	q.Lock()
	defer q.Unlock()
	q.data = append(q.data, *data)
	q.index++
	if q.index > len(q.data) {
		q.index = 0
	}
	if q.length < len(q.data) {
		q.length++
	}
	q.refreshExp(expDuration)
}

func (q *Cache) All() (values []models.Event) {
	q.RLock()
	defer q.RUnlock()
	q.refreshExp(expDuration)
	return append(q.data[q.index:len(q.data)], q.data[:q.index]...)
}

func (q *Cache) Top() *models.Event {
	q.RLock()
	defer q.RUnlock()
	q.refreshExp(expDuration)
	return &q.data[q.index]
}

func (q *Cache) Limit(length int) (values []models.Event) {
	q.RLock()
	defer q.RUnlock()
	if length < len(q.data) {
		return q.data
	}
	var result []models.Event
	if q.index > length {
		result = q.data[q.index-length : q.index]
	} else {
		result = append(q.data[len(q.data)-length+q.index:len(q.data)], q.data[:q.index]...)
	}
	q.refreshExp(expDuration)
	return result
}
