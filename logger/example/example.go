package main

import (
	"os"

	"github.com/PrathamSkilltelligent/pmgo/logger"
)

func main() {
	isProd := false
	log := logger.NewLogger(os.Stdout, "felix", isProd)
	defer log.Close()

	log.Info("This is an info message", "key", "value", "count", 42)
	log.Error("This is an error message", "error", "example error")
	log.Warn("This is a warning message", "warning", "example warning")
	log.Debug("This is a debug message", "debug", "example debug")
}
