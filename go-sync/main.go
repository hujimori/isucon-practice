package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	s := time.Now()
	fmt.Printf("[Start]\t\t[%v]\n", getMinutesAndSeconds())
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go test(i, &wg)
	}
	fmt

	fmt.Printf("CPUのコア数: %v\n", runtime.NumCPU())
	fmt.Printf("Goroutineの数: %v\n", runtime.NumGoroutine())

	for i := 0; i < 5; i++ {
		time.Sleep(time.Second * 1)
		fmt.Printf("[%v]メインでも処理中...\n", i)
	}
	wg.Wait()
	fmt.Printf("[Done]\t\t[%v]\n", getMinutesAndSeconds())

	e := time.Now()
	fmt.Printf("処理秒数: %v\n", e.Sub(s).Round(time.Second))
}

func getMinutesAndSeconds() string {
	return time.Now().Format("04:05")
}

func test(n int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[Start]\t%v\t[%v]\n", n, getMinutesAndSeconds())
	time.Sleep(time.Second * 15)
	fmt.Printf("[Done]\t%v\t[%v]\n", n, getMinutesAndSeconds())
}
