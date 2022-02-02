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
	capacity   int
	queue      List
	itemsByKey map[Key]*ListItem
	sync.RWMutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	if listItem, ok := c.itemsByKey[key]; ok {
		listItem.Value = cacheItem{key: key, value: value}
		c.queue.MoveToFront(listItem)
		return true
	}

	listItem := c.queue.PushFront(cacheItem{key: key, value: value})
	c.itemsByKey[key] = listItem
	if c.queue.Len() > c.capacity {
		itemToRemove := c.queue.Back()

		c.queue.Remove(itemToRemove)

		delete(c.itemsByKey, itemToRemove.Value.(cacheItem).key)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	if listItem, ok := c.itemsByKey[key]; ok {
		c.queue.MoveToFront(listItem)
		return listItem.Value.(cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.queue = NewList()
	c.itemsByKey = make(map[Key]*ListItem, c.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:   capacity,
		queue:      NewList(),
		itemsByKey: make(map[Key]*ListItem, capacity),
	}
}
