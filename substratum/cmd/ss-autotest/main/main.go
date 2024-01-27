package main

import (
	"fmt"
	"os"

	ssautotest "github.com/sgtcodfish/substratum/cmd/ss-autotest"
	"github.com/sgtcodfish/substratum/cmd/util"
)

func main() {
	ctx := util.Context()

	err := ssautotest.Invoke(ctx, os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
