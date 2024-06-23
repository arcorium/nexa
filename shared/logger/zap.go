package logger

import (
  "fmt"
  "go.uber.org/zap"
)

func NewZapLogger(isDev bool, options ...zap.Option) (Logger, error) {
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
    logger: logger,
  }, nil
}

type ZapLogger struct {
  logger *zap.Logger
}

func (z *ZapLogger) Debugf(format string, args ...interface{}) {
  z.logger.Debug(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Infof(format string, args ...interface{}) {
  z.logger.Info(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Warnf(format string, args ...interface{}) {
  z.logger.Warn(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Fatalf(format string, args ...interface{}) {
  z.logger.Fatal(fmt.Sprintf(format, args...))
}

func (z *ZapLogger) Debug(s string) {
  z.logger.Debug(s)
}

func (z *ZapLogger) Info(s string) {
  z.logger.Info(s)
}

func (z *ZapLogger) Warn(s string) {
  z.logger.Warn(s)
}

func (z *ZapLogger) Fatal(s string) {
  z.logger.Fatal(s)
}
