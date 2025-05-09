// Package fast 提供高性能的JSON序列化和反序列化功能。
package fast

import (
	"sync"
)

// shardCount 是分片锁的数量，用于减少锁竞争。
const shardCount = 32

// fragmentCache 用于存储大型JSON对象的片段缓存。
type fragmentCache struct {
	shards [shardCount]*cacheShard
}

// cacheShard 是缓存分片。
type cacheShard struct {
	mu    sync.RWMutex
	cache map[string]interface{}
}

// globalFragmentCache 是全局片段缓存实例。
var globalFragmentCache = newFragmentCache()

// newFragmentCache 创建新的片段缓存。
func newFragmentCache() *fragmentCache {
	fc := &fragmentCache{}
	for i := 0; i < shardCount; i++ {
		fc.shards[i] = &cacheShard{
			cache: make(map[string]interface{}, 64), // 预分配合理大小。
		}
	}
	return fc
}

// getShard 获取key对应的分片。
func (fc *fragmentCache) getShard(key string) *cacheShard {
	// 使用简单的哈希算法选择分片。
	hash := fnvHash(key)
	return fc.shards[hash%shardCount]
}

// fnvHash 实现FNV-1a哈希算法。
func fnvHash(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= prime32
	}
	return hash
}

// Set 存储片段。
func (fc *fragmentCache) Set(key string, value interface{}) {
	shard := fc.getShard(key)
	shard.mu.Lock()
	shard.cache[key] = value
	shard.mu.Unlock()
}

// Get 获取片段。
func (fc *fragmentCache) Get(key string) (interface{}, bool) {
	shard := fc.getShard(key)
	shard.mu.RLock()
	value, ok := shard.cache[key]
	shard.mu.RUnlock()
	return value, ok
}

// Delete 删除片段。
func (fc *fragmentCache) Delete(key string) {
	shard := fc.getShard(key)
	shard.mu.Lock()
	delete(shard.cache, key)
	shard.mu.Unlock()
}

// Clear 清空缓存。
func (fc *fragmentCache) Clear() {
	for i := 0; i < shardCount; i++ {
		shard := fc.shards[i]
		shard.mu.Lock()
		shard.cache = make(map[string]interface{}, 64) // 预分配合理大小。
		shard.mu.Unlock()
	}
}

// CacheFragment 缓存JSON片段。
func CacheFragment(key string, value interface{}) {
	globalFragmentCache.Set(key, value)
}

// GetCachedFragment 获取缓存的JSON片段。
func GetCachedFragment(key string) (interface{}, bool) {
	return globalFragmentCache.Get(key)
}

// ClearFragmentCache 清空片段缓存。
func ClearFragmentCache() {
	globalFragmentCache.Clear()
}
