package logger

import (
  "fmt"
  "go.uber.org/zap"
)

func NewZapLogger(isDev bool, options ...zap.Option) (ILogger, error) {
  var logger *zap.Logger
  var err error
  if isDev {
    logger, err = zap.NewDevelopment(options...)
    if err != nil {
      return nil, err
    }
  } else {
    logger, err = zap.NewProduction(options...)
    if err != nil {
      return nil, err
    }
  }

  return &ZapLogger{
    Internal: logger,
  }, nil
}

type ZapLogger struct {
  Internal *zap.Logger
}

func (z *ZapLogger) Debugf(format string, args ...interface{}) {
  z.Internal.Debug(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Infof(format string, args ...interface{}) {
  z.Internal.Info(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Warnf(format string, args ...interface{}) {
  z.Internal.Warn(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
  z.Internal.Fatal(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Debug(s string) {
  z.Internal.Debug(s)
}

func (z *ZapLogger) Info(s string) {
  z.Internal.Info(s)
}

func (z *ZapLogger) Warn(s string) {
  z.Internal.Warn(s)
}

func (z *ZapLogger) Fatal(s string) {
  z.Internal.Fatal(s)
}
