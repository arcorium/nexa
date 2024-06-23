package interceptor

import (
  "context"
  "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
  "go.uber.org/zap"
)

func ZapLogger(logger *zap.Logger) logging.Logger {
  return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
    fs := make([]zap.Field, 0, len(fields))

    for i := 0; i < len(fields); i += 2 {
      key := fields[i]
      value := fields[i+1]

      switch v := value.(type) {
      case string:
        fs = append(fs, zap.String(key.(string), v))
      case int:
        fs = append(fs, zap.Int(key.(string), v))
      case bool:
        fs = append(fs, zap.Bool(key.(string), v))
      default:
        fs = append(fs, zap.Any(key.(string), value))
      }
    }

    l := logger.WithOptions(zap.AddCallerSkip(1)).With(fs...)

    switch lvl {
    case logging.LevelDebug:
      l.Debug(msg)
    case logging.LevelInfo:
      l.Info(msg)
    case logging.LevelWarn:
      l.Warn(msg)
    case logging.LevelError:
      l.Error(msg)
    }
  })
}
