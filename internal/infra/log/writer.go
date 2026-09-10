package log

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type writerSet struct {
	writer  io.Writer
	closers []io.Closer
}

func newWriter(cfg Config) (writerSet, error) {
	var writers []io.Writer
	var closers []io.Closer

	if cfg.Output.Stdout {
		writers = append(writers, terminalWriter(cfg, os.Stdout))
	}
	if cfg.Output.Stderr {
		writers = append(writers, terminalWriter(cfg, os.Stderr))
	}

	if cfg.Output.File != nil {
		f := &lumberjack.Logger{
			Filename: cfg.Output.File.Filename, MaxSize: cfg.Output.File.MaxSize,
			MaxBackups: cfg.Output.File.MaxBackups, MaxAge: cfg.Output.File.MaxAge,
			Compress: cfg.Output.File.Compress, LocalTime: cfg.Output.File.LocalTime,
		}
		writers = append(writers, f)
		closers = append(closers, f)
	}

	switch len(writers) {
	case 0:
		return writerSet{}, fmt.Errorf("no log writer configured")
	case 1:
		return writerSet{writer: writers[0], closers: closers}, nil
	default:
		return writerSet{writer: zerolog.MultiLevelWriter(writers...), closers: closers}, nil
	}
}

func terminalWriter(cfg Config, out io.Writer) io.Writer {
	if cfg.Format == FormatConsole {
		return zerolog.ConsoleWriter{Out: out, NoColor: !cfg.Color, TimeFormat: cfg.TimeFormat}
	}
	return out
}
