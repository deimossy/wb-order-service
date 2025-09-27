package inmemory

import (
	"sync"
	"time"
)

type Cache[T any] interface {
	Set(key string, value T, ttl time.Duration)
	Get(key string) (T bool)
	Delete(key string)
	Clear()
	Len() int
}

type cacheItem[T any] struct {
	value      T
	expiration int64
}

type MemoryCache[T any] struct {
	data map[string]cacheItem[T]
	mu   sync.RWMutex

	capacity int
}

func NewMemoryCache[T any]() *MemoryCache[T] {
	c := &MemoryCache[T]{
		data:     make(map[string]cacheItem[T]),
		capacity: 0,
	}

	go c.cleanupLoop()

	return c
}

func (c *MemoryCache[T]) Set(key string, value T, ttl time.Duration) {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	} else {
		exp = int64(1<<63 - 1)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity > 0 && len(c.data) >= c.capacity {
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}

	c.data[key] = cacheItem[T]{
		value: value,
		expiration: exp,
	}
}

func (c *MemoryCache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	it, ok := c.data[key]
	c.mu.RUnlock()
	if !ok {
		var zero T
		return zero, false
	}

	if time.Now().UnixNano() > it.expiration {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		var zero T
		return zero, false
	}

	return it.value, true
}

func (c *MemoryCache[T]) Delete(key string) {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
}

func (c *MemoryCache[T]) Clear() {
	c.mu.Lock()
	c.data = make(map[string]cacheItem[T])
	c.mu.Unlock()
}

func (c *MemoryCache[T]) Len() int {
	c.mu.RLock()
	l := len(c.data)
	c.mu.RUnlock()
	return l
}

func (c *MemoryCache[T]) cleanupLoop() {
	t := time.NewTicker(1 * time.Minute)
	defer t.Stop()
	for range t.C {
		now := time.Now().UnixNano()
		c.mu.Lock()
		for k, v := range c.data {
			if now > v.expiration {
				delete(c.data, k)
			}
		}

		c.mu.Unlock()
	}
}
