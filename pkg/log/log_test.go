package log

import (
	"context"
	"testing"
)

func TestLog(t *testing.T) {
	Init(&Options{
		OutputPaths:       []string{"stderr"},
		Level:             "DEBUG",
		Format:            "json",
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       false,
		Development:       false,
	})
	Debug("debug test")
	L(context.Background()).Info("=====================")
	Info("info test")
	Error("error test")
}
