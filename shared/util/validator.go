package util

import (
	"context"
	"github.com/go-playground/validator/v10"
	"nexa/shared/status"
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

func ValidateStruct[T any](ctx context.Context, strct *T) status.Object {
	err := GetValidator().StructCtx(ctx, strct)
	if err != nil {
		return status.Error(status.FIELD_VALIDATION_ERROR, err)
	}
	return status.SuccessInternal()
}
