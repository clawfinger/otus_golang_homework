package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.RWMutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	if listItem, ok := c.items[key]; ok {
		listItem.Value = value
		c.queue.MoveToFront(listItem)
		return true
	} else {
		c.items[key] = c.queue.PushFront(value)

		if c.queue.Len() >= c.capacity {
			c.queue.Remove(c.queue.Back())
		}
		return false
	}

}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	if listItem, ok := c.items[key]; ok {
		c.queue.MoveToFront(listItem)
		return listItem.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
