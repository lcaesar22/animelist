package validator

import (
	"regexp"
	"slices"
)

var EmailRX = regexp.MustCompile("^[-zA*Z0-9.!#$%&'*+/?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Error map[string]string
}

// New is a helper which creates a validation instance
func New() *Validator {
	return &Validator{Error: make(map[string]string)}
}

// Valid returns true if the map doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.Error) == 0
}

// AddError adds an error message to the map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Error[key]; exists {
		v.Error[key] = message
	}
}

// Check adds an error message to the map only if the validation check is not okay
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Generic function which returns true if all values in a slice are unique
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches returns true if a string value matches a specific regex pattern
func Matches(value string, regex *regexp.Regexp) bool {
	return regex.MatchString(value)
}

// Generic value which returns true if all the values in a slice are unique
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(uniqueValues) == len(values)
}
