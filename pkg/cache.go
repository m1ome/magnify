package pkg

import (
	"sync"
	"time"
)

type (
	Cache struct {
		c   sync.Map
		exp time.Duration
	}

	CacheValue struct {
		exp time.Time
		v   any
	}
)

func NewCache(exp time.Duration) *Cache {
	return &Cache{
		c:   sync.Map{},
		exp: exp,
	}
}

func (c *Cache) Store(key string, v any) {
	c.c.Store(key, CacheValue{exp: time.Now().Add(c.exp), v: v})
}

func (c *Cache) Load(key string) (any, bool) {
	v, ok := c.c.Load(key)
	if !ok {
		return nil, false
	}

	return v.(CacheValue).v, true
}

func (c *Cache) Copy() map[string]any {
	m := make(map[string]any)
	c.c.Range(func(key, value any) bool {
		m[key.(string)] = value.(CacheValue).v
		return true
	})

	return m
}

func (c *Cache) Cleanup() {
	c.c.Range(func(key, value any) bool {
		if time.Now().After(value.(CacheValue).exp) {
			c.c.Delete(key)
		}

		return true
	})
}
