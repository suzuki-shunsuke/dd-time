package main

import (
	"os"

	"github.com/suzuki-shunsuke/dd-time/pkg/cli"
)

func main() {
	if code := cli.Core(); code != 0 {
		os.Exit(code)
	}
}
