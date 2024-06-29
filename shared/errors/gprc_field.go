package errors

import (
  "errors"
  "fmt"
  "github.com/go-playground/validator/v10"
  "google.golang.org/genproto/googleapis/rpc/errdetails"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "log"
)

type FieldError struct {
}

func GrpcFieldErrors(details ...*errdetails.BadRequest_FieldViolation) error {
  badreq := &errdetails.BadRequest{}
  badreq.FieldViolations = append(badreq.FieldViolations, details...)

  st := status.New(codes.InvalidArgument, "Invalid Argument")
  stats, err := st.WithDetails(badreq)
  if err != nil {
    log.Fatalln("Fatal error: ", err)
  }

  return stats.Err()
}

func GrpcFieldValidationErrors(verr validator.ValidationErrors) error {
  var result []*errdetails.BadRequest_FieldViolation
  for _, val := range verr {
    result = append(result, &errdetails.BadRequest_FieldViolation{
      Field:       val.StructField(),
      Description: val.Error(),
    })
  }

  return GrpcFieldErrors(result...)
}

func GrpcFieldIndexedErrors(field string, ierr IndicesError) error {
  if errors.Is(ierr.Errs[0], ErrEmptySlice) {
    return GrpcFieldErrors(&errdetails.BadRequest_FieldViolation{
      Field:       field,
      Description: fmt.Sprintf("%s should not be empty", field),
    })
  }

  var result []*errdetails.BadRequest_FieldViolation
  for key, val := range ierr.Errs {
    result = append(result, &errdetails.BadRequest_FieldViolation{
      Field:       fmt.Sprintf("%s[%d]", field, key),
      Description: val.Err.Error(),
    })
  }

  return GrpcFieldErrors(result...)
}
