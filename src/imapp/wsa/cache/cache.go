package cache

import (
	"github.com/dgraph-io/ristretto"
	"sync"
	"time"
)

type WACache struct {
	*ristretto.Cache
}

var (
	once  sync.Once
	cache *WACache
)

const (
	DefaultExpireTime = 10 * time.Second
)

func Cache() *WACache {
	once.Do(func() {
		config := ristretto.Config{
			NumCounters: 2,      // number of keys to track frequency of (10M).
			MaxCost:     2,      // maximum cost of cache (1GB).
			BufferItems: 2 << 5, // number of keys per Get buffer.
		}

		c, _ := ristretto.NewCache(&config)
		cache = &WACache{
			Cache: c,
		}
	})
	return cache
}
