package clap

import "go.uber.org/zap"

// TODO: make no-op log
var log, _ = zap.NewDevelopment()

// SetLogger adds a global log for the clap package.
func SetLogger(l *zap.Logger) {
	if l == nil {
		return
	}
	log = l
}
