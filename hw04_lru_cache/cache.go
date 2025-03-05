package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    *List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value any
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value any) bool {
	if item, exists := c.items[key]; exists {
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)
		return true
	}

	if c.capacity == 0 {
		return false
	}

	if c.queue.Len() == c.capacity {
		back := c.queue.Back()
		c.queue.Remove(back)
		delete(c.items, back.Value.(*cacheItem).key)
	}

	newCacheItem := &cacheItem{
		key:   key,
		value: value,
	}
	newListItem := c.queue.PushFront(newCacheItem)
	c.items[key] = newListItem
	return false
}

func (c *lruCache) Get(key Key) (any, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
