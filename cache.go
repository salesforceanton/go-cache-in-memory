package cache

import (
	"crypto/sha1"
	"fmt"
	"time"
)

const (
	cryptoKey = "dsafsdg23dds"
)

type Cache struct {
	scope map[string]interface{}
}

func (c *Cache) Set(key string, value interface{}) {
	c.scope[key] = value
}
func (c *Cache) SetWithLifetime(key string, value interface{}, lifetime time.Duration) {
	c.scope[key] = value
	go c.RunCleaner(key, lifetime)
}
func (c *Cache) Get(key string) interface{} {
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
	time.Sleep(lifetime)
	c.Delete(key)
}

func NewCache() *Cache {
	return &Cache{
		scope: make(map[string]interface{}),
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
