package clap

import (
	"go.uber.org/zap"
	"os"
)

var log = zap.NewNop()

// GetLogger returns the dev logger for the clap package.
func GetLogger() *zap.Logger {
	return log
}

// SetLogger sets the dev logger for the clap package.
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
