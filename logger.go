package clap

import (
	"go.uber.org/zap"
	"os"
)

// TODO: make no-op log
var log = zap.NewNop()

// SetLogger adds a global log for the clap package.
func SetLogger(l *zap.Logger) {
	if l == nil {
		return
	}
	log = l
}

func init() {
	_, exists := os.LookupEnv("CLAP_DEBUG")

	if !exists {
		return
	}

	l, _ := zap.NewDevelopment()
	SetLogger(l)
}