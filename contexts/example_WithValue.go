package main

import (
	"context"
	"fmt"
)

type favContextKey string

func CheckValue(ctx context.Context, k favContextKey) {
	var value string
	if v := ctx.Value(k); v != nil {
		u, ok := v.(string)
		if !ok {
			fmt.Println("Not authorized")
			return
		}
		value = u
		fmt.Println("found value:", value)
		return
	}

	fmt.Println("key not found:", k)
}

func main() {
	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")

	CheckValue(ctx, k)
	CheckValue(ctx, favContextKey("color"))

	dl, ok := ctx.Deadline()
	fmt.Printf("deadline : %v\nok : %v\n", dl, ok)

}
