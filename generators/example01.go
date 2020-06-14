package main

import "fmt"

func combination(ch chan []string, str, arr []string, n, r, index int) {

	if r == 0 {
		ch <- append(str[:0:0], str...)
		return
	} else if n == 0 || n < r {
		return
	} else {

		str = append(str, arr[index])
		combination(ch, str, arr, n-1, r-1, index+1)
		str = str[:len(str)-1]
		combination(ch, str, arr, n-1, r, index+1)
	}
}

func Combinations(arr []string, m int) <-chan []string {
	ch := make(chan []string)
	var str []string

	go func() {
		defer close(ch)
		combination(ch, str, arr, len(arr), m, 0)
	}()

	return ch
}

func main() {
	for c := range Combinations([]string{"사과", "배", "복숭아", "포도", "귤"}, 3) {
		fmt.Println(c)
	}
}
