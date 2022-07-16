package main

import "fmt"

func main() {
	ss := make([]int, 20)
	for i := 0; i < 20; i++ {
		ss[i] = i
	}
	fmt.Print(ss)
	fmt.Print(ss[:10])
}
