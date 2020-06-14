package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	d := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	defer cancel()

	select {
	case <-time.After(101 * time.Millisecond):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // prints "context deadline exceeded"
	}

	dl, ok := ctx.Deadline()
	fmt.Printf("deadline : %v\nok : %v\n", dl, ok)

}
