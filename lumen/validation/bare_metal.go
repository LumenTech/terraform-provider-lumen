package validation

import (
	"regexp"
	"strings"
)

const hostnameLengthMessage = "A hostname should be between 1 and 253 characters comprised of letters, numbers, hyphens, and periods."
const hostnameRegexMessage = "Each element of the hostname, separated by a period, should be at most 63 characters and should not begin with a hyphen."

var hostnameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]{0,62}([.][a-zA-Z0-9][a-zA-Z0-9-]{0,62})*$")

type CustomValidationError struct {
	messages []string
}

func (v CustomValidationError) Error() string {
	return strings.Join(v.messages, "\n")
}

func ValidateBareMetalServerName(name string) error {
	var validationErrors []string
	if len(name) < 1 || len(name) > 253 {
		validationErrors = append(validationErrors, hostnameLengthMessage)
	}

	matched := hostnameRegex.Match([]byte(name))
	if !matched {
		validationErrors = append(validationErrors, hostnameRegexMessage)
	}

	if len(validationErrors) > 0 {
		return CustomValidationError{
			messages: validationErrors,
		}
	}
	return nil
}
