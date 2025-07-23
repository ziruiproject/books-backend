package validator

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	idTranslation "github.com/go-playground/validator/v10/translations/id"
	"github.com/rs/zerolog/log"
	ivalidator "starter/internal/core/validator"
)

type Validator struct {
	Instance *validator.Validate
	Uni      *ut.UniversalTranslator
	message  map[string][]ivalidator.ValidationError
}

func (v *Validator) InitTranslation() {

	v.message = make(map[string][]ivalidator.ValidationError)
	v.message["en"] = []ivalidator.ValidationError{
		{
			Message: "cannot be empty",
			Field:   "required",
		},
	}
	v.message["id"] = []ivalidator.ValidationError{
		{
			Message: "tidak boleh kosong",
			Field:   "required",
		},
	}

	english, found := v.Uni.GetTranslator("en")
	if !found {
		log.Error().Msgf(`Translation for "en" not found`)
	}
	if err := enTranslation.RegisterDefaultTranslations(v.Instance, english); err != nil {
		log.Err(err).Msgf(`Failed to load translation for "en" not found`)
	}

	indonesian, found := v.Uni.GetTranslator("id")
	if !found {
		log.Error().Msgf(`Translation for "id" not found, defaulting to "en"`)
	} else {
		if err := idTranslation.RegisterDefaultTranslations(v.Instance, indonesian); err != nil {
			log.Err(err).Msgf(`Failed to load translation for "id" not found`)
		}
	}

	err := v.RegisterTranslations(english, "en")
	if err != nil {
		log.Err(err).Msgf(`Failed to register "en" translation`)
		return
	}

	err = v.RegisterTranslations(indonesian, "id")
	if err != nil {
		log.Err(err).Msgf(`Failed to register "id" translation`)
		return
	}
}

func (v *Validator) RegisterTranslations(trans ut.Translator, locale string) error {
	for _, message := range v.message[locale] {
		err := v.Instance.RegisterTranslation(message.Field, trans, func(ut ut.Translator) error {
			return ut.Add(message.Field, fmt.Sprintf("{0} %s", message.Message), true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(message.Field, fe.Field())
			return t
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) GetTranslator(locale string) (ut.Translator, error) {
	trans, found := v.Uni.GetTranslator(locale)
	if !found {
		return nil, fmt.Errorf("translation for %s not found", locale)
	}
	return trans, nil
}

func (v *Validator) ValidateStruct(item interface{}) []ivalidator.ValidationError {
	var errors []ivalidator.ValidationError

	err := v.Instance.Struct(item)
	if err == nil {
		return nil
	}

	trans, errTrans := v.GetTranslator("en")
	if errTrans != nil {
		trans, _ = v.GetTranslator("en")
	}

	for _, failed := range err.(validator.ValidationErrors) {
		errors = append(errors, ivalidator.ValidationError{
			Message: failed.Translate(trans),
			Field:   failed.Field(),
		})
	}

	return errors
}

func NewValidator() ivalidator.Validator {
	english := en.New()
	indonesian := id.New()
	uni := ut.New(english, english, indonesian)

	validation := &Validator{
		Instance: validator.New(validator.WithRequiredStructEnabled()),
		Uni:      uni,
		message:  make(map[string][]ivalidator.ValidationError),
	}
	validation.InitTranslation()

	return validation
}
