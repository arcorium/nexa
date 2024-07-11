package status

import (
  "errors"
  spanUtil "github.com/arcorium/nexa/shared/util/span"
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

func NewWithMessage(code Code, msg string) Object {
  return Object{
    Codes: code,
    Error: errors.New(msg),
  }
}

type Object struct {
  Codes Code
  Error error
}

func (s *Object) Message() string {
  if s.Error != nil {
    return s.Error.Error()
  }
  return ""
}

func (s *Object) IsError() bool {
  return s.Codes > DELETED
}

// ToGRPCError convert status into GRPC error, it will return nil when the Code is not an error code
// which will be handled by grpc status package
func (s *Object) ToGRPCError() error {
  return status.Error(MapGRPCCode(s.Codes), s.Message())
}

func (s *Object) ToGRPCErrorWithSpan(span trace.Span) error {
  var err error
  if s.IsError() {
    code := MapGRPCCode(s.Codes)
    err = status.Error(code, s.Message())
    spanUtil.RecordError(err, span)
  } else {
    span.SetStatus(codes.Ok, s.Message())
  }
  return err
}
