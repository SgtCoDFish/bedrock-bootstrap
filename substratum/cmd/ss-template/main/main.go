package main

import (
	"fmt"
	"os"

	sstemplate "github.com/sgtcodfish/substratum/cmd/ss-template"
	"github.com/sgtcodfish/substratum/cmd/util"
)

func main() {
	ctx := util.Context()

	err := sstemplate.Invoke(ctx, os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse flags / setup: %v\n", err)
		os.Exit(1)
	}
}
