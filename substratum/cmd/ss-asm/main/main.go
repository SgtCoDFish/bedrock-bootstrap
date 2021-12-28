package main

import (
	"context"
	"fmt"
	"os"

	ssasm "github.com/sgtcodfish/substratum/cmd/ss-asm"
)

func main() {
	ctx := context.Background()

	err := ssasm.Invoke(ctx, os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse flags / setup: %v\n", err)
		os.Exit(1)
	}
}
