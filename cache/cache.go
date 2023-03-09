package cache

/*
	基于 github.com/patrickmn/go-cache 再次封装的内存缓存
	用于一般场景下的内存缓存

*/

import (
	"fmt"
	cache "github.com/patrickmn/go-cache"
	"time"
)

type MemCache struct {
	C           *cache.Cache
	DefaultTime time.Duration // 默认过期时间
	CleanTime   time.Duration // 每 CleanTime时间清除一次过期项目
}

// NewCache 创建一个默认过期时间为 defaultTime 实际的缓存，每 cleanTime时间清除一次过期项目
func NewCache(defaultTime, cleanTime time.Duration) *MemCache {
	return &MemCache{
		C:           cache.New(defaultTime, cleanTime),
		DefaultTime: defaultTime,
		CleanTime:   cleanTime,
	}
}

// Set 写入数据
func (mem *MemCache) Set(key string, value interface{}) error {
	if len(key) < 1 {
		return fmt.Errorf("key cannot null")
	}
	mem.C.Set(key, value, mem.DefaultTime)
	return nil
}

// SetExp 写入数据并设置过期时间
func (mem *MemCache) SetExp(key string, value interface{}, expiration time.Duration) error {
	if len(key) < 1 {
		return fmt.Errorf("key cannot null")
	}
	mem.C.Set(key, value, expiration)
	return nil
}

// Clear 清除缓存
func (mem *MemCache) Clear() {
	mem.C.Flush()
}

// Get 获取缓存
func (mem *MemCache) Get(key string) (interface{}, bool) {
	return mem.C.Get(key)
}

// Delete 删除
func (mem *MemCache) Delete(key string) {
	mem.C.Delete(key)
}

// GetAll 获取全部
func (mem *MemCache) GetAll() map[string]cache.Item {
	return mem.C.Items()
}

// Save 持久化到文件
func (mem *MemCache) Save(filePath string) error {
	return mem.C.SaveFile(filePath)
}

// Load 加载持久化
func (mem *MemCache) Load(filePath string) error {
	return mem.C.LoadFile(filePath)
}
