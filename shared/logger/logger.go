package logger

import "nexa/shared/wrapper"

var Default = wrapper.Must(NewZapLogger(true))

func SetGlobal(logger Logger) {
  Default = logger
}

type Logger interface {
  Debugf(format string, args ...any)
  Infof(format string, args ...any)
  Warnf(format string, args ...any)
  Fatalf(format string, args ...any)
  Debug(string)
  Info(string)
  Warn(string)
  Fatal(string)
}

func Debugf(format string, args ...any) {
  Default.Debugf(format, args...)
}

func Infof(format string, args ...any) {
  Default.Infof(format, args...)
}

func Warnf(format string, args ...any) {
  Default.Warnf(format, args...)
}

func Fatalf(format string, args ...any) {
  Default.Fatalf(format, args...)
}

func Debug(s string) {
  Default.Debug(s)
}

func Info(s string) {
  Default.Infof(s)
}

func Warn(s string) {
  Default.Warn(s)
}

func Fatal(s string) {
  Default.Fatalf(s)
}
