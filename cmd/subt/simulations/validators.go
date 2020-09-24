package simulations

import (
	"gopkg.in/go-playground/validator.v9"
)

// InstallSimulationCustomValidators extends validator.v9 with custom validation
// functions and meta tags for simulations.
func InstallSimulationCustomValidators(validate *validator.Validate) {
	_ = validate.RegisterValidation("isruletype", isValidRuleType)
}

// CustomRuleTypes contains the list of available rule types
var CustomRuleTypes = []CustomRuleType{
	MaxSubmissions,
}

// IsValidRuleType checks
func isValidRuleType(fl validator.FieldLevel) bool {
	value := CustomRuleType(fl.Field().String())
	for _, r := range CustomRuleTypes {
		if value == r {
			return true
		}
	}
	return false
}
