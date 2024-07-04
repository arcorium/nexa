package span

import (
  "go.opentelemetry.io/otel/codes"
  "go.opentelemetry.io/otel/trace"
)

// RecordError Set span status into codes.Error and record error for trace.Span.
// This function doesn't check the error, use it only when the error is not nil
func RecordError(err error, spans trace.Span) {
  spans.SetStatus(codes.Error, err.Error())
  spans.RecordError(err)
}
