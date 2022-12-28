package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"unicode"

	"github.com/dimfeld/httptreemux/v5"
	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"

	// v9 DOCs : https://pkg.go.dev/gopkg.in/go-playground/validator.v9?utm_source=godoc
	validator "gopkg.in/go-playground/validator.v9"
	//validator "github.com/go-playground/validator/v10"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	//en_translations Nope. v10 breaks all this.
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator *ut.UniversalTranslator

func init() {

	// Instantiate the validator for use.
	validate = validator.New()

	// Instantiate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, lang)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Custom validation function for keywords and categories.
	validate.RegisterValidation("keywords", keywords)
}

const KEYWORDS_MAX_LENGTH = 50

// keywords returns true if string is only alphanum,
// and optionally whitespaces and hyphens, and within KEYWORDS_MAX_LENGTH.
func keywords(fl validator.FieldLevel) bool {

	// if strings.ContainsAny(fl.Field().String(), "\"~`!@#$%^&*()_+=[]{}\\|;:'\"/?>.<,") {
	// 	return false
	// }

	//flag := false
	count := 0
	safe := func(r rune) rune {
		count = count + 1
		if count > KEYWORDS_MAX_LENGTH {
			return -1
		}
		switch {
		case r > unicode.MaxASCII:
			return -1
		case unicode.IsLetter(r):
			return r
		case unicode.IsNumber(r):
			return r
		case rune(' ') == r:
			// if flag {
			// 	return -1
			// }
			// flag = true
			return r
		case rune('-') == r:
			// if flag {
			// 	return -1
			// }
			// flag = true
			return r
		default:
			return -1
		}
	}
	for strings.Map(safe, fl.Field().String()) != fl.Field().String() {
		return false
	}
	return true
}

// Params returns the web call parameters from the request.
func Params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}

// Decode reads the body of an HTTP request (r) looking for a JSON document.
// The body is decoded (unmarshalled) into the provided pointer (ptr).
// If pointing to a struct, then validator operates against its field tags.
func Decode(r *http.Request, ptr interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(ptr); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(ptr); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls that of the error messages.
		// Add "Accept-Language" header code here to support multiple languages.
		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}
