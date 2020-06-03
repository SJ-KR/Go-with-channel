package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func readword(ch chan string, file *os.File) {
	defer file.Close()

	r := bufio.NewReader(file)
	for {
		word, err := r.ReadString(' ')
		ch <- word
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return
		}
	}
}

func allwords(filename string) <-chan string {
	ch := make(chan string)
	file, err := os.OpenFile(filename, os.O_RDONLY, os.FileMode(644))

	if err != nil {
		log.Println(err)
	}

	go func() {
		defer close(ch)
		readword(ch, file)
	}()

	return ch
}

func main() {
	filepath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	filepath += "/sample.txt"
	fmt.Println(filepath)
	for w := range allwords(filepath) {
		fmt.Println(w)

	}
}
