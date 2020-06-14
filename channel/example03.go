package main

import "fmt"

func hanoigen(ch chan [2]string, n int, from, to, by string) {
	if n == 0 {
		return
	}
	hanoigen(ch, n-1, from, by, to)
	ch <- [2]string{from, to}
	hanoigen(ch, n-1, by, to, from)

}
func Hanoi(n int, from, to, by string) chan [2]string {
	ch := make(chan [2]string)

	go func() {
		defer close(ch)
		hanoigen(ch, n, from, to, by)
	}()

	return ch
}

func main() {
	for move := range Hanoi(3, "A", "B", "C") {
		fmt.Println(move[0], "->", move[1])
	}
}
