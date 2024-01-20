package cache

import (
	"container/list"
	"errors"
	"sync"
)

var (
	ErrInvalidCapacity = errors.New("invalid capacity")
	ErrKeyNotFound     = errors.New("key not found")
)

type Cache struct {
	capacity int
	cache    *list.List
	elements map[string]*list.Element
	mutex    sync.RWMutex
}

func New(capacity int) (*Cache, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	return &Cache{
		capacity: capacity,
		cache:    list.New(),
		elements: make(map[string]*list.Element),
	}, nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[string(key)]; ok {
		value := elem.Value.(Item).Value
		c.cache.MoveToFront(elem)
		return value, nil
	}
	return nil, ErrKeyNotFound
}

func (c *Cache) Set(key, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[string(key)]; ok {
		c.cache.MoveToFront(elem)
		elem.Value = Item{
			Key:   key,
			Value: val,
		}
	} else {
		if c.cache.Len() == c.capacity {
			index := c.cache.Back().Value.(Item).Key
			delete(c.elements, string(index))
			c.cache.Remove(c.cache.Back())
		}
	}

	item := &list.Element{Value: Item{
		Key:   key,
		Value: val,
	}}

	i := c.cache.PushFront(item)
	c.elements[string(key)] = i

}

func (c *Cache) Remove(key []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[string(key)]; ok {
		delete(c.elements, string(key))
		c.cache.Remove(elem)
		return nil
	}
	return ErrKeyNotFound
}
