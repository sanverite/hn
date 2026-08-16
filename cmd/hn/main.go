package main

import (
	"os"

	"github.com/sanverite/hn/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
