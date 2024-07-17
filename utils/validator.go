package utils

import "github.com/go-playground/validator/v10"

type ParamError struct {
	Param   string `json:"param"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func validatorTagMapper(tag string) string {
	switch tag {
	case "gte":
		return "O numero deve ser positivo"
	case "required":
		return "Este campo é obrigatório"
	}
	return tag
}

func Validator(s interface{}) []ParamError {
	validate := validator.New()
	err := validate.Struct(s)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		out := make([]ParamError, len(validationErrors))
		for i, e := range validationErrors {
			out[i] = ParamError{
				e.Field(),
				validatorTagMapper(e.Tag()),
				e.Type().Kind().String(),
			}
		}
		return out
	}
	return nil
}
