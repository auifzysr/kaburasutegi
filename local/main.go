package main

import (
	"log/slog"
	"os"

	"github.com/auifzysr/kaburasutegi/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
