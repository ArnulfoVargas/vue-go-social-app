package validator

import (
	"Server/internal/dto"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	return &Validator{v: validator.New()}
}

func (val *Validator) Validate(s any) []dto.Error {
	var errs []dto.Error
	err := val.v.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, dto.Error{
				Field: e.Field(),
				Tag:   e.Tag(),
			})
		}
	}
	return errs
}
