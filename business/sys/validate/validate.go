package validate

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"reflect"
	"strings"
)

// validate holds the settings and caches for validating request structs.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

// init validator and translator.
func init() {
	// init a validator
	validate = validator.New()

	// create a translator for english
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// register english error messages
	_ = entranslations.RegisterDefaultTranslations(validate, translator)

	// use json tag names for errors instead of struct names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
}

// Check validates the provided model against its declared tags.
func Check(val any) error {
	if err := validate.Struct(val); err != nil {
		// use a type assertion to get the real error message
		ferrs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// parse errors
		var fields FieldErrors
		for _, ferr := range ferrs {
			field := FieldError{
				Field: ferr.Field(),
				Error: ferr.Translate(translator),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}

// GenerateId generates a unique id for entities.
func GenerateId() string {
	return uuid.NewString()
}

// CheckId validates that the format of an id is valid.
func CheckId(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidId
	}

	return nil
}

// todo: CheckEmail validates that the format of an email is valid.
