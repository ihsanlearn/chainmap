package main

import (
	"os"

	"github.com/ihsanlearn/chainmap/logger"
	"github.com/ihsanlearn/chainmap/options"
	"github.com/ihsanlearn/chainmap/pkg/runner"
)

func main() {
	opts := options.ParseOptions()
	r := runner.New(opts)

	if err := r.CheckDependencies(); err != nil {
		logger.Error("System check failed: %s", err)
		os.Exit(1)
	}

	r.Run()
}
