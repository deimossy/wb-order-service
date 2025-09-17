package inmemory

import "time"

// Cache - интерфейс кэша
type Cache[T any] interface {
	Set(key string, value T, ttl time.Duration)
	Get(key string) (T, bool)
	Delete(key string)
	Clear()
}