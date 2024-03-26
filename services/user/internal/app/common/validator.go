package common

import (
	"github.com/go-playground/validator/v10"
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
