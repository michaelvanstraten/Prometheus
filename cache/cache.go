package cache

import (
	"sort"
	"sync"
)

type CacheItem struct {
	data []byte
	rank int64
} 

type Cache struct {
	mux sync.Mutex
	MemoryLimit int64
	db map[string]CacheItem
}

func New() *Cache {
	var NewCache = &Cache{}
	NewCache.db = make(map[string]CacheItem)
	NewCache.MemoryLimit = 64000000
	NewCache.mux = sync.Mutex{}
	return NewCache
}

func (c *Cache) FlushDB() {
	c.db = make(map[string]CacheItem)
}

func (c *Cache) RemoveLeastImportentItem() {
	var dbSlice = make([]struct {
		Value CacheItem 
		Key string
	}, 0, len(c.db))
	for key, data := range c.db {
		dbSlice = append(dbSlice, struct{Value CacheItem; Key string}{
			Value: data,
			Key: key,
		})
	}
	sort.Slice(dbSlice, func(i, j int) bool {
		if dbSlice[i].Value.rank != dbSlice[j].Value.rank {
			return dbSlice[i].Value.rank > dbSlice[j].Value.rank 
		} else {
			return false
		}
	})
	c.mux.Lock()
	var item = dbSlice[len(dbSlice)-1]
	delete(c.db, item.Key)
	c.MemoryLimit += int64(len(item.Value.data))
	c.mux.Unlock()
}

func (c *Cache) Set(Key string, Data []byte) {
	c.mux.Lock()
	if item, ok := c.db[Key]; ok {
		item.data = Data
		c.db[Key] = item
		c.mux.Unlock()
		return
	}
	c.mux.Unlock()
	var sizeOfData int64 = int64(len(Data))
	for sizeOfData > c.MemoryLimit {
		if len(c.db) == 0 {
			return
		}
		c.RemoveLeastImportentItem()
	}
	c.mux.Lock()
	c.db[Key] = CacheItem {
		data : Data,
		rank: sizeOfData,
	}
	c.MemoryLimit -= sizeOfData
	c.mux.Unlock()
}

func (c *Cache) Get(Key string) []byte {
	if item, ok := c.db[Key]; ok {
		return item.data
	} else {
		return nil
	}
}