package main

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.Mutex
	items map[int]int
}

func NewCache() *Cache {
	m := make(map[int]int)
	c := &Cache{
		items: m,
	}
	return c
}

func (c *Cache) Set(key int, value int) {
	c.mu.Lock()
	c.items[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key int) int {
	c.mu.Lock()
	v, ok := c.items[key]
	c.mu.Unlock()

	// キャッシュにはデータがある場合、即座にリターン
	if ok {
		return v
	}

	// 無い場合は、重い処理（データベースへのアクセス）が行われる
	v = HeavyGet(key)

	c.Set(key, v)

	return v
}

func HeavyGet(key int) int {
	time.Sleep(time.Second)
	return key * 2
}

var mcache *Cache

func main() {
	mcache = NewCache()

	go func() {
		for i := 0; i < 100; i++ {
			log.Printf("goroutine %d\n", mcache.Get(i))
		}
	}()

	for i := 0; i < 100; i++ {
		log.Println(mcache.Get(i))
	}
}
