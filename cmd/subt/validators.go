package main

import (
	"encoding/json"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

// This module adds custom validators used by validator.v9

const (
	// Matches alphanum chars plus underscore, dash and spaces (\t\n\f\r )
	alphaNumSpaceUnderscoreDash = "^[\\w\\-\\s]+$"
)

var (
	alphaNumSpaceUnderscoreDashRegex = regexp.MustCompile(alphaNumSpaceUnderscoreDash)
	allowedExpFeatures               = []string{"gzweb"}
)

var blacklist []string

// InstallCustomValidators extends validator.v9 with custom validation functions
// and meta tags for fields.
func InstallCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation("noforwardslash", notIncludeForwardSlash)
	validate.RegisterValidation("alphanumspace", isAlphanumSpace)
	loadBlacklist()
	validate.RegisterValidation("notinblacklist", notInBlacklist)
	validate.RegisterValidation("expfeatures", isExpFeatures)
	validate.RegisterValidation("lenEqFieldLen", lenEqFieldLen)
}

func loadBlacklist() {
	data, err := ioutil.ReadFile("validators_owners_blacklist.json")
	if err != nil {
		log.Fatal("Couldn't read blacklist file", err)
		return
	}
	err = json.Unmarshal(data, &blacklist)
	if err != nil {
		log.Fatal("Couldn't unmarshal blacklist", err)
		return
	}
}

// notInBlacklist is the validation function for validating if the current
// field's value is not listed in the blacklist of owner names.
// From: https://github.com/marteinn/The-Big-Username-Blacklist
func notInBlacklist(fl validator.FieldLevel) bool {
	return !includeString(fl.Field().String(), blacklist)
}

func includeString(val string, list []string) bool {
	for _, s := range list {
		if s == val {
			return true
		}
	}
	return false
}

// isAlphanumSpace is the validation function for validating if the current
// field's value is a valid alphanumeric value that also accepts dashes,
// underscores and spaces.
func isAlphanumSpace(fl validator.FieldLevel) bool {
	return alphaNumSpaceUnderscoreDashRegex.MatchString(fl.Field().String())
}

// notIncludeForwardSlash is a function that validates the field value does not
// include forward slashes (/).
func notIncludeForwardSlash(fl validator.FieldLevel) bool {
	return !strings.Contains(fl.Field().String(), "/")
}

// isExpFeatures is a function that validates if the field's value is a comma
// separated list of words, and that each word belongs to the
// expFeatures whitelist.
// If the input is empty, the validation will be OK too.
func isExpFeatures(fl validator.FieldLevel) bool {
	features := ign.StrToSlice(fl.Field().String())
	if len(features) == 0 {
		return true
	}
	for _, f := range features {
		if !includeString(f, allowedExpFeatures) {
			return false
		}
	}
	return true
}

// lenEqFieldLen validates if the current field's len() is equal to the len() of
// the field specified by the param's value.
// TODO: move this to ign-fuelserver's validator.go or to a common place (eg. ign-go).
func lenEqFieldLen(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, ok := fl.GetStructFieldOK()
	if !ok || currentKind != kind {
		return false
	}
	return currentField.Len() == field.Len()
}
