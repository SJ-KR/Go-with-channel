package main

import "fmt"

func gray_reverse(ch chan []int, buf []int, n, index int) {
	if n == index {
		ch <- append(buf[:0:0], buf...)
		return
	}
	buf[index] = 1
	gray(ch, buf, n, index+1)
	buf[index] = 0
	gray_reverse(ch, buf, n, index+1)
}
func gray(ch chan []int, buf []int, n, index int) {
	if n == index {
		ch <- append(buf[:0:0], buf...)
		return
	}
	buf[index] = 0
	gray(ch, buf, n, index+1)
	buf[index] = 1
	gray_reverse(ch, buf, n, index+1)
}
func GrayBinaryGenerator(n int) <-chan []int {
	ch := make(chan []int)
	buf := make([]int, n)

	go func() {
		defer close(ch)
		gray(ch, buf, n, 0)
	}()

	return ch
}

func main() {
	for g := range GrayBinaryGenerator(4) {
		fmt.Println(g)
	}
}
