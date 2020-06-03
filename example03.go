package main

import "fmt"

func hanoigen(ch chan []string, n int, from, to, by string) {
	if n == 0 {
		return
	}

	hanoigen(ch, n-1, from, by, to)
	ch <- []string{from, to}
	hanoigen(ch, n-1, by, to, from)

}
func Hanoi(n int, from, to, by string) chan []string {
	ch := make(chan []string)

	go func() {
		defer close(ch)
		hanoigen(ch, n, from, by, to)
	}()

	return ch
}

func main() {
	for move := range Hanoi(3, "A", "B", "C") {
		fmt.Println(move[0], "->", move[1])
	}
}
