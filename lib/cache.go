package lib

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type CacheItem struct {
	Key        string
	Value      interface{}
	Expiration int64 // 过期时间，单位：秒
}

type Cache struct {
	items    map[string]*CacheItem
	mu       sync.RWMutex
	stop     chan struct{}
	interval time.Duration // 清理过期 item 的时间间隔
}

// NewCache 返回一个新的 Cache 对象
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		items:    make(map[string]*CacheItem),
		stop:     make(chan struct{}),
		interval: interval,
	}

	go cache.startCleanup()

	return cache
}

// Set 向缓存中添加一个新的 item
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl).Unix()
	c.items[key] = &CacheItem{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}
}

// Get 从缓存中获取指定 key 的 item，如果不存在则返回错误
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.items[key]; ok {
		if time.Now().Unix() > item.Expiration {
			delete(c.items, key)
			return nil, errors.New("item has expired")
		}
		return item.Value, nil
	}

	return nil, fmt.Errorf("item not found for key: %s", key)
}

// Delete 从缓存中删除指定 key 的 item
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// startCleanup 定期清理过期的 item
func (c *Cache) startCleanup() {
	ticker := time.NewTicker(c.interval)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, item := range c.items {
				if time.Now().Unix() > item.Expiration {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}

// Stop 停止缓存清理
func (c *Cache) Stop() {
	c.stop <- struct{}{}
}
