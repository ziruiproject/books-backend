package validator

import ut "github.com/go-playground/universal-translator"

type Validator interface {
	InitTranslation()
	GetTranslator(locale string) (ut.Translator, error)
	ValidateStruct(item interface{}) []ValidationError
	RegisterTranslations(trans ut.Translator, locale string) error
}
