package main

import (
	"github.com/liwanggui/dnscli-go/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
