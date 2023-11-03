package validation

import (
	"regexp"
	"strings"
	"unicode"
)

const hostnameLengthMessage = "A hostname should be between 1 and 253 characters comprised of letters, numbers, hyphens, and periods."
const hostnameRegexMessage = "Each element of the hostname, separated by a period, should be at most 63 characters and should not begin with a hyphen."

var hostnameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]{0,62}([.][a-zA-Z0-9][a-zA-Z0-9-]{0,62})*$")
var alphaNumberDashesAndUnderscoresRegex = regexp.MustCompile("^[A-Za-z0-9_-]+$")

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

func ValidateBareMetalUsername(username string) error {
	var validationErrors []string
	if len(username) == 0 {
		validationErrors = append(validationErrors, "Please provide a valid username.")
	}

	if strings.ToLower(username) == "root" {
		validationErrors = append(validationErrors, "The username 'root' is not permitted.")
	}

	if !alphaNumberDashesAndUnderscoresRegex.Match([]byte(username)) {
		validationErrors = append(validationErrors, "The username can only contain alphanumeric characters, dashes, and underscores.")
	}

	if len(validationErrors) > 0 {
		return CustomValidationError{
			messages: validationErrors,
		}
	}
	return nil
}

var passwordValidationError = CustomValidationError{
	messages: []string{
		`Please provide a password that conforms to the provided rules.
Must be at least 9 characters
* uppercase letters
* lowercase letters
* numbers
`,
	},
}

var passwordMustIncludeTests = []func(rune) bool{
	unicode.IsUpper,
	unicode.IsLower,
	unicode.IsDigit,
}

var passwordMustNotIncludeTests = []func(rune) bool{
	unicode.IsSymbol,
	unicode.IsPunct,
	unicode.IsSpace,
}

func ValidateBareMetalPassword(password string) error {
	if len(password) < 9 {
		return passwordValidationError
	}

	for _, test := range passwordMustIncludeTests {
		found := false
		for _, r := range password {
			if test(r) {
				found = true
			}
		}
		if !found {
			return passwordValidationError
		}
	}

	for _, test := range passwordMustNotIncludeTests {
		found := false
		for _, r := range password {
			if test(r) {
				found = true
			}
		}
		if found {
			return passwordValidationError
		}
	}

	return nil
}
