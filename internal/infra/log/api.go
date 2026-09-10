package log

import "github.com/rs/zerolog"

func Trace() *zerolog.Event {
	l := L()
	return l.Trace()
}

func Debug() *zerolog.Event {
	l := L()
	return l.Debug()
}

func Info() *zerolog.Event {
	l := L()
	return l.Info()
}

func Warn() *zerolog.Event {
	l := L()
	return l.Warn()
}

func Error() *zerolog.Event {
	l := L()
	return l.Error()
}

func Fatal() *zerolog.Event {
	l := L()
	return l.Fatal()
}

func Panic() *zerolog.Event {
	l := L()
	return l.Panic()
}
