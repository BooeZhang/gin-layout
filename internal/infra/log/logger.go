package log

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
	closers []io.Closer
}

// New 新建一个 log 对象
func New(cfg Config) (*Logger, error) {
	cfg, err := cfg.validate()
	if err != nil {
		return nil, err
	}

	ws, err := newWriter(cfg)
	if err != nil {
		return nil, err
	}

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		for _, closer := range ws.closers {
			_ = closer.Close()
		}
		return nil, err
	}

	ctx := zerolog.New(ws.writer).Level(level).With().Timestamp()
	if cfg.Service != "" {
		ctx = ctx.Str("service", cfg.Service)
	}
	if cfg.Environment != "" {
		ctx = ctx.Str("environment", cfg.Environment)
	}
	if cfg.Version != "" {
		ctx = ctx.Str("version", cfg.Version)
	}

	l := &Logger{Logger: ctx.Logger().Hook(contextHook{}), closers: ws.closers}
	if cfg.Caller {
		l.Logger = l.With().Caller().Logger()
	}

	return l, nil
}

func (l *Logger) Close() error {
	if l == nil {
		return nil
	}
	var firstErr error
	for _, closer := range l.closers {
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
