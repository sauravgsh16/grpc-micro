package main

import (
	"fmt"
	"os"

	"github.com/sauravgsh16/api-grpc/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
