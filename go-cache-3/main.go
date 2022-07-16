package main

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.Mutex
	items map[int]int
	time  int64
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
	c.time = time.Now().Unix()
	c.mu.Unlock()
}

func (c *Cache) Get(key int) int {
	c.mu.Lock()

	v, ok := c.items[key]
	c.mu.Unlock()

	now := time.Now().Unix()

	if ok {
		if now-c.time >= 30 {
			log.Printf("cache is expired %d - %d = %d\n", now, c.time, now-c.time)

		} else {
			return v
		}
	}

	go func() {
		log.Println("Heavy Get...")
		v := HeavyGet(key)
		c.Set(key, v)
	}()

	return 100
}

func HeavyGet(key int) int {
	time.Sleep(time.Second)
	return key * 2
}

func main() {
	mcache := NewCache()
	go func() {
		for i := 0; i < 10000; i++ {
			log.Println(mcache.Get(i))
			log.Printf("goroutine\n")
		}
	}()

	for i := 0; i < 100; i++ {

		time.Sleep(time.Second * 10)

		log.Printf("Get %d\n", mcache.Get(i))
	}
}
