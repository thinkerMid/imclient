package memoryCacheService

import (
	"github.com/dgraph-io/ristretto"
	"sync"
	"time"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/application/data_storage/cache/constant"
)

var _ containerInterface.IMemoryCacheService = &MemoryCache{}

var cacheTTL = time.Minute * 5
var globalCache *ristretto.Cache
var cacheInit sync.Once

// 经测试如果每个账号各自持有一个cache 内部协程开销以及初始化开销是比较多 不划算
// 改回原有的公用缓存池与各自持有的方式相比有以下缺点：
// 1.账号应用释放的时候，所缓存的内存数据不能及时释放，需要等待ttl
// 2.数据负载不均衡，比如一个账号可以存储2000个数据，但是公用情况下不可能一样，多的会很多少的会很少
// 3.公用数据池满载的情况下，非热点数据不容易一直缓存
func cache() *ristretto.Cache {
	cacheInit.Do(func() {
		config := ristretto.Config{
			// 每个账号存2048个数据 最多2000个账号
			NumCounters: 2048 * 2048,
			// 每个账号0.5 MB的容量 最多2000个账号
			MaxCost: (2 << 18) * 2048,
			// 64 B
			BufferItems: 2 << 5,
		}
		globalCache, _ = ristretto.NewCache(&config)
	})

	return globalCache
}

// MemoryCache 数据缓存服务
type MemoryCache struct {
	containerInterface.BaseService

	accountLoginData *cacheConstant.AccountLoginData
}

// Init .
func (m *MemoryCache) Init() {
	m.accountLoginData = &cacheConstant.AccountLoginData{ABExposureKey: []byte{}}
}

// Cache .
func (m *MemoryCache) Cache(key, value interface{}) {
	cache().SetWithTTL(key, value, 1, cacheTTL)
}

// CacheTTL .
func (m *MemoryCache) CacheTTL(key, value interface{}, ttl time.Duration) {
	cache().SetWithTTL(key, value, 1, ttl)
}

// UnCache .
func (m *MemoryCache) UnCache(key interface{}) {
	cache().Del(key)
}

// FindInCache .
func (m *MemoryCache) FindInCache(key interface{}) (interface{}, bool) {
	return cache().Get(key)
}

// AccountLoginData .
func (m *MemoryCache) AccountLoginData() *cacheConstant.AccountLoginData {
	return m.accountLoginData
}

// OnApplicationExit .
func (m *MemoryCache) OnApplicationExit() {}
