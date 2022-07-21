package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	go insertBackground()

	for i := 0; i < 1000; i++ {
		time.Sleep(time.Second)
		insert(i, i*2, i*3)
	}

	for {

	}
}

var insertQuech = make(chan struct {
	params []interface{}
}, 1000)

var once sync.Once

type InsertQue struct {
	query  []string
	params []interface{}
}

func insertBackground() {
	var insertQue InsertQue

	// query := make([]string, 0, 1000)
	// params := make([]interface{}, 0, 1000)

	insertQue.query = make([]string, 0, 1000)
	insertQue.params = make([]interface{}, 0, 1000)

	ticker := time.NewTicker(10 * time.Second)
	once.Do(func() {
		fmt.Println("Start!")
	})

	for {
		select {
		case <-ticker.C:
			insertFunc := func() {
				fmt.Printf("query:%v params:%v\n", insertQue.query, insertQue.params)
				// for i := 0; i < len(query); i++ {

				// }

			}
			insertFunc()
			// insertQue.query
			insertQue.query = make([]string, 0, 1000)
			insertQue.params = make([]interface{}, 0, 1000)
			// params = make([]interface{}, 0, 1000)

		case v := <-insertQuech:
			insertQue.query = append(insertQue.query, "(?, ?, ?)")
			insertQue.params = append(insertQue.params, v.params...)

			fmt.Printf("キュー追加：%v\n", v)
		}

	}

}

func insert(p1 int, p2 int, p3 int) {
	insertQuech <- struct {
		params []interface{}
	}{params: []interface{}{p1, p2, p3}}
}
