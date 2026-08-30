package main

import (
	"fmt"
	"os"

	"github.com/Mort4lis/booking-tennis-court/internal/cli"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	if err := cli.NewRootCmd(version).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
