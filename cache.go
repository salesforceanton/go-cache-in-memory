package cache

import (
	"crypto/sha1"
	"fmt"
	"sync"
	"time"
)

const (
	cryptoKey = "dsafsdg23dds"
)

type Cache struct {
	scope map[string]interface{}
	mu    *sync.RWMutex
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	c.scope[key] = value
	c.mu.Unlock()
}
func (c *Cache) SetWithLifetime(key string, value interface{}, lifetime time.Duration) {
	c.mu.Lock()
	c.scope[key] = value
	c.mu.Unlock()
	go c.RunCleaner(key, lifetime)
}
func (c *Cache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.scope[key]
	if !ok {
		return nil
	}
	return value
}
func (c *Cache) Delete(key string) {
	delete(c.scope, key)
}
func (c *Cache) RunCleaner(key string, lifetime time.Duration) {
	for {
		select {
		case <-time.After(lifetime):
			c.Delete(key)
			return

		default:
		}
	}
}

func NewCache() *Cache {
	return &Cache{
		scope: make(map[string]interface{}),
		mu:    new(sync.RWMutex),
	}
}

// TODO: Draft logic need to improve
func CreateCacheableFunc(targetFunc func(...interface{}) interface{}) func(...interface{}) interface{} {
	cache := NewCache()
	return func(params ...interface{}) interface{} {
		key := cache.generateKeyHash(params...)

		value := cache.Get(key)
		if value == nil {
			value = targetFunc(params...)
			cache.Set(key, value)
			fmt.Printf("Put result to cache")
		} else {
			fmt.Printf("Retain result from cache")
		}
		return value
	}
}
func (c *Cache) generateKeyHash(params ...interface{}) string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprint(params...)))

	return fmt.Sprintf("%x", hash.Sum([]byte(cryptoKey)))
}
