package main

import (
	"container/list"
	"sync"
)


type Item struct {
	Key []byte
	Value []byte
}

type Cache struct{
	capacity int
	cache *list.List
	elements map[string]*list.Element
	mutex sync.RWMutex
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache: list.New(),
		elements: make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key []byte) []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.elements[string(key)]; ok{
		value:=elem.Value.(*list.Element).Value.(Item).Value
		c.cache.MoveToFront(elem)
		return value
	}
	return nil
}


func main() {
	
}