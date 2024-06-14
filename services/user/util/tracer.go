package util

import (
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
  "nexa/services/user/constant"
)

func GetTracer(options ...trace.TracerOption) trace.Tracer {
  t := otel.GetTracerProvider()
  return t.Tracer(constant.TracerName, options...)
}
