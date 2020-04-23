package tools

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// ParseJSONStruct parses the given Request body as JSON with the given Struct type.
func ParseJSONStruct(s interface{}, r *http.Request) *ign.ErrMsg {
	if err := json.NewDecoder(r.Body).Decode(s); err != nil {
		return ign.NewErrorMessageWithBase(ign.ErrorUnmarshalJSON, err)
	}
	return nil
}

// ParseFormStruct parses the given Request body as Multi-part form using the given formDecoder and the given Struct type.
func ParseFormStruct(s interface{}, r *http.Request, formDecoder *form.Decoder) *ign.ErrMsg {
	if errs := formDecoder.Decode(s, r.Form); errs != nil {
		return ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs,
			getDecodeErrorsExtraInfo(errs))
	}
	return nil
}

// ValidateStruct uses the given validator to ensure that it given struct is valid.
// It's usually used after parsing a struct with ParseJSONStruct and ParseFormStruct.
func ValidateStruct(s interface{}, validator *validator.Validate) *ign.ErrMsg {
	if errs := validator.Struct(s); errs != nil {
		return ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs,
			getValidationErrorsExtraInfo(errs))
	}
	return nil
}

func getDecodeErrorsExtraInfo(err error) []string {
	errs := err.(form.DecodeErrors)
	extra := make([]string, 0, len(errs))
	for field, er := range errs {
		extra = append(extra, fmt.Sprintf("Field: %s. %v", field, er.Error()))
	}
	return extra
}

func getValidationErrorsExtraInfo(err error) []string {
	validationErrors := err.(validator.ValidationErrors)
	extra := make([]string, 0, len(validationErrors))
	for _, fe := range validationErrors {
		extra = append(extra, fmt.Sprintf("%s:%v", fe.StructField(), fe.Value()))
	}
	return extra
}
