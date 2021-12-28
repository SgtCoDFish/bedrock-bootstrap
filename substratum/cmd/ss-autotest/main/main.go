package main

import (
	"context"
	"fmt"
	"os"

	ssautotest "github.com/sgtcodfish/substratum/cmd/ss-autotest"
)

func main() {
	ctx := context.Background()

	err := ssautotest.Invoke(ctx, os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse flags / setup: %v\n", err)
		os.Exit(1)
	}
}
