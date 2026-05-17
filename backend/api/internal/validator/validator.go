package validator

import (
	"Server/internal/domain"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	return &Validator{v: validator.New()}
}

func (val *Validator) Validate(s any) []domain.Error {
	var errs []domain.Error
	err := val.v.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, domain.Error{
				Field: e.Field(),
				Tag:   e.Tag(),
			})
		}
	}
	return errs
}
