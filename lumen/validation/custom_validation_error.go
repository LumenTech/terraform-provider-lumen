package validation

import "strings"

type CustomValidationError struct {
	messages []string
}

func (v CustomValidationError) Error() string {
	return strings.Join(v.messages, "\n")
}
