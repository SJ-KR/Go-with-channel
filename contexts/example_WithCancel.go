package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func gen(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		defer close(dst)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("called <-ctx.Done()")
				return
			case dst <- n:
				time.Sleep(time.Millisecond * 10)
				n++
			}
		}
	}()
	return dst
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		fmt.Println("started first")
		for n := range gen(ctx) {
			fmt.Println(n)
			if n == 2 {
				fmt.Println("ended first")
				break
			}
		}
		wg.Done()
	}()

	cancel()

	go func() {
		fmt.Println("started second")
		for n := range gen(ctx) {
			fmt.Println(n)
			if n == 3 {
				fmt.Println("ended second")
				break
			}
		}
		wg.Done()
	}()

	wg.Wait()

	dl, ok := ctx.Deadline()
	fmt.Println("----------------------------------")
	fmt.Printf("deadline : %v\nok : %v\n", dl, ok)
	fmt.Printf("number of go routine : %d\n", runtime.NumGoroutine())
}
