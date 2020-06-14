package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

type result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func walkFiles(ctx context.Context, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)

		errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-ctx.Done():
				return errors.New("walk canceled")
			}
			return nil
		})
	}()
	return paths, errc
}
func digester(ctx context.Context, paths <-chan string, c chan<- result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case c <- result{path, md5.Sum(data), err}:
		case <-ctx.Done():
			fmt.Println("digester canceled")
			return
		}
	}
}

func MD5All(ctx context.Context, root string) (map[string][md5.Size]byte, error) {
	paths, errc := walkFiles(ctx, root)

	c := make(chan result)
	var wg sync.WaitGroup

	const numDigesters = 20
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(ctx, paths, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	m := make(map[string][md5.Size]byte)
	for r := range c {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}

	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func main() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Millisecond))
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	m, err := MD5All(ctx, os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}

	dl, ok := ctx.Deadline()
	fmt.Println("----------------------------------")
	fmt.Printf("deadline : %v\nok : %v\n", dl, ok)
	fmt.Printf("number of go routine : %d\n", runtime.NumGoroutine())
}
