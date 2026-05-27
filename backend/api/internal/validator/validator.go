package validator

import (
	"Server/internal/shared"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	return &Validator{v: validator.New()}
}

func (val *Validator) Validate(s any) []shared.Error {
	var errs []shared.Error
	err := val.v.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, shared.Error{
				Field: e.Field(),
				Tag:   e.Tag(),
			})
		}
	}
	return errs
}
