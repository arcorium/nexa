package status

import (
  "go.opentelemetry.io/otel/codes"
  "go.opentelemetry.io/otel/trace"
  "google.golang.org/grpc/status"
)

func Error(code Code, err error) Object {
  return New(code, err)
}

//func FieldError() Object {
//
//}

func ErrorC(code Code) Object {
  return New(code, nil)
}

func New(code Code, err error) Object {
  return Object{
    Codes: code,
    Error: err,
  }
}

type Object struct {
  Codes Code
  Error error
}

func (s *Object) IsError() bool {
  return s.Error != nil || s.Codes > DELETED
}

// ToGRPCError convert status into GRPC error, it will return nil when the Code is not an error code
// which will be handled by grpc status package
func (s *Object) ToGRPCError() error {
  return status.Error(MapGRPCCode(s.Codes), s.Error.Error())
}

func (s *Object) ToGRPCErrorWithSpan(span trace.Span) error {
  span.RecordError(s.Error)
  err := status.Error(MapGRPCCode(s.Codes), s.Error.Error())
  if err != nil {
    span.SetStatus(codes.Error, err.Error())
  }
  return err
}
