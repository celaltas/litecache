package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	capacity int
	cache    *list.List
	elements map[string]*list.Element
	mutex    sync.RWMutex
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache:    list.New(),
		elements: make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key []byte) []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[string(key)]; ok {
		value := elem.Value.(*list.Element).Value.(Item).Value
		c.cache.MoveToFront(elem)
		return value
	}
	return nil
}

func (c *Cache) Set(key, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[string(key)]; ok {
		c.cache.MoveToFront(elem)
		elem.Value.(*list.Element).Value = Item{
			Key:   key,
			Value: val,
		}
	} else {
		if c.cache.Len() == c.capacity {
			index := c.cache.Back().Value.(*list.Element).Value.(Item).Key
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

func (c *Cache) Remove(key []byte) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[string(key)]; ok {
		delete(c.elements, string(key))
		c.cache.Remove(elem)
	}
}
