package cache

import (
	"sync"
)

type Servicer[T interface{}] interface {
	GetOne(key string) (*T, error)
}

type Cache[T interface{}] struct {
	service  Servicer[T]
	dataLock sync.RWMutex
	data     map[string]*T
}

func (c *Cache[T]) get(id string) (*T, bool) {
	c.dataLock.RLock()
	defer c.dataLock.RUnlock()

	user, ok := c.data[id]
	return user, ok
}

func (c *Cache[T]) set(id string, value *T) {
	c.dataLock.Lock()
	defer c.dataLock.Unlock()

	c.data[id] = value
}

func (c *Cache[T]) GetOne(key string) (*T, error) {
	value, ok := c.get(key)
	if !ok {
		value, err := c.service.GetOne(key)
		if err != nil {
			return nil, err
		}

		c.set(key, value)
		return value, nil
	}

	return value, nil
}

func NewCache[T interface{}](service Servicer[T]) *Cache[T] {
	return &Cache[T]{
		service: service,
		data:    map[string]*T{},
	}
}
