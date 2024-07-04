package util

import (
  "context"
  "github.com/go-playground/validator/v10"
  "golang.org/x/exp/constraints"
  sharedErr "nexa/shared/errors"
  "nexa/shared/types"
  "sync"
)

var validatorInstanceOnce sync.Once
var validatorInstance *validator.Validate

func GetValidator() *validator.Validate {
  validatorInstanceOnce.Do(func() {
    validatorInstance = validator.New()
  })
  return validatorInstance
}

func ValidateStructCtx[T any](ctx context.Context, strct *T) error {
  err := GetValidator().StructCtx(ctx, strct)
  if err == nil {
    return nil
  }

  verr, ok := err.(validator.ValidationErrors)
  if !ok {
    return err
  }
  return sharedErr.GrpcFieldValidationErrors(verr)
}

func ValidateStruct[T any](strct *T) error {
  err := GetValidator().Struct(strct)
  if err == nil {
    return nil
  }

  verr, ok := err.(validator.ValidationErrors)
  if !ok {
    return err
  }
  return sharedErr.GrpcFieldValidationErrors(verr)
}

func StringEmptyValidates(fields ...types.Field[string]) sharedErr.EmptyFieldError {
  var errs []sharedErr.FieldError

  for _, field := range fields {
    if len(field.Val) == 0 {
      errs = append(errs, sharedErr.NewFieldError(field.Name, sharedErr.ErrFieldEmpty))
    }
  }
  return sharedErr.EmptyFieldError{Errs: errs}
}

func ZeroIntegerValidates[T constraints.Integer](fields ...types.Field[T]) sharedErr.EmptyFieldError {
  var errs []sharedErr.FieldError

  for _, field := range fields {
    if field.Val == 0 {
      errs = append(errs, sharedErr.NewFieldError(field.Name, sharedErr.ErrZeroEmpty))
    }
  }
  return sharedErr.EmptyFieldError{Errs: errs}
}

func ZeroFloatValidates[T constraints.Float](fields ...types.Field[T]) sharedErr.EmptyFieldError {
  var errs []sharedErr.FieldError

  for _, field := range fields {
    if field.Val == 0.0 {
      errs = append(errs, sharedErr.NewFieldError(field.Name, sharedErr.ErrZeroEmpty))
    }
  }
  return sharedErr.EmptyFieldError{Errs: errs}
}
