package logic

import "sync"

type Cache[K comparable, V any] struct {
	lock sync.RWMutex
	data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		data: map[K]V{},
	}
}

func (c *Cache[K, V]) Put(key K, value V) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	result := false
	if _, ok := c.data[key]; !ok {
		result = true
		c.data[key] = value
	}
	return result
}

func (c *Cache[K, V]) Replace(key K, value V) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	previous, ok := c.data[key]
	c.data[key] = value
	return previous, ok
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	return value, ok
}
