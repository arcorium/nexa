package util

import (
  "context"
  "github.com/go-playground/validator/v10"
  sharedErr "nexa/shared/errors"
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
