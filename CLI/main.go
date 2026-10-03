package main

import (
	"fmt"
	"os"

	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/app"
	"github.com/Amonnium/Calculate-Pi-Scripts/cli/internal/tui"
)

func main() {
	application, err := app.New(os.Args[1:])
	if err == nil {
		err = tui.Run(application)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
