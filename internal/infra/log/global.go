package log

import (
	"errors"
	"sync/atomic"

	"github.com/rs/zerolog"
)

var ErrAlreadyInitialized = errors.New("log already initialized")

var global atomic.Pointer[Logger]

// Init 初始化全局 log
func Init(cfg Config) error {
	if global.Load() != nil {
		return ErrAlreadyInitialized
	}

	cfg, err := cfg.validate()
	if err != nil {
		return err
	}

	zerolog.TimeFieldFormat = cfg.TimeFormat

	l, err := New(cfg)
	if err != nil {
		return err
	}

	if !global.CompareAndSwap(nil, l) {
		_ = l.Close()
		return ErrAlreadyInitialized
	}

	zerolog.DefaultContextLogger = &l.Logger

	return nil
}

func L() zerolog.Logger {
	if l := global.Load(); l != nil {
		return l.Logger
	}
	return zerolog.Nop()
}

func Close() error {
	l := global.Swap(nil)
	if l == nil {
		return nil
	}
	zerolog.DefaultContextLogger = nil
	return l.Close()
}
