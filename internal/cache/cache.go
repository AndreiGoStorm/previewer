package cache

import (
	"container/list"
	"fmt"
	"path"
	"strings"
	"sync"

	"previewer/internal/config"
	"previewer/internal/service"
)

type Cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    *list.List
	items    map[string]*list.Element
	storage  *service.Storage
	sync.Mutex
}

type Item struct {
	key   string
	value interface{}
}

func New(conf config.Cache, storage *service.Storage) Cache {
	lru := &lruCache{
		capacity: conf.Capacity,
		queue:    list.New(),
		items:    make(map[string]*list.Element, conf.Capacity),
		storage:  storage,
	}

	err := lru.warmingCache()
	if err != nil {
		panic(err)
	}

	return lru
}

func (lru *lruCache) warmingCache() error {
	names, err := lru.storage.ReadDirNames()
	if err != nil {
		return err
	}

	var ext string
	for _, name := range names {
		ext = path.Ext(name)
		lru.Set(strings.TrimRight(name, ext), ext)
	}

	return nil
}

func (lru *lruCache) Set(key string, value interface{}) bool {
	defer lru.Mutex.Unlock()
	lru.Mutex.Lock()
	newItem := &Item{key, value}
	item, ok := lru.items[key]
	if ok {
		item.Value = newItem
		lru.moveToFront(item, key)
	} else {
		lru.dequeueOldestItem()
		lru.items[key] = lru.queue.PushFront(newItem)
	}

	return ok
}

func (lru *lruCache) Get(key string) (interface{}, bool) {
	defer lru.Mutex.Unlock()
	lru.Mutex.Lock()
	if item, ok := lru.items[key]; ok {
		lru.queue.MoveToFront(item)
		currentItem := item.Value.(*Item)
		return currentItem.value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	defer lru.Mutex.Unlock()
	lru.Mutex.Lock()
	clear(lru.items)
	lru.queue = list.New()
}

func (lru *lruCache) moveToFront(item *list.Element, key string) {
	lru.queue.MoveToFront(item)
	lru.items[key] = lru.queue.Front()
}

func (lru *lruCache) dequeueOldestItem() {
	if lru.queue.Len() >= lru.capacity {
		last := lru.queue.Back()
		item := last.Value.(*Item)
		delete(lru.items, item.key)
		lru.queue.Remove(last)

		filename := fmt.Sprintf("%s%v", item.key, item.value)
		_ = lru.storage.DeleteFile(filename)
	}
}
