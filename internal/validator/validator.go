package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Validator is a struct for implementing validation logic and holding validation error strings.
type Validator struct {
	NonFieldErrors []string          // validation errors not related to a specific form field e.g. login failure
	FieldErrors    map[string]string // validation errors related to a specific form field, where form field name is key
}

// Valid returns true if FieldErrors does not contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// CheckField adds an error message to the FieldErrors map only if a
// validation check is not ok.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// AddFieldError sets key and message into FieldErrors, but only if
// key does not exists. If key exists, this operation is a no-op.
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// AddNonFieldError appends message to NonFieldErrors.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// NotBlank returns true if the value is not an empty string after
// removing blank spaces.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinChars returns true if value contains at least n characters
func MinChars(value string, n int) bool {
	fmt.Println(utf8.RuneCountInString(value))
	return utf8.RuneCountInString(value) >= n
}

// MaxChars returns true if value contains less than or equals to n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt returns true if value is within the slice of permittedValues.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Matches returns true if a value matches the provided regex ex.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
