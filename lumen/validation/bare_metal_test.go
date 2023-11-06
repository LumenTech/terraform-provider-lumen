package validation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestValidateBareMetalServerName_ValidInputs(t *testing.T) {
	validInputs := []string{
		"mybaremetal",
		"abcdefghijk-01235",
		"67890l-mnopqrstuvwxyz",
		"multiple.domain.levels",
		"m",
		strings.Repeat("a", 63),
		strings.Repeat(strings.Repeat("a", 49)+".", 5) + "com",
	}

	for _, input := range validInputs {
		err := ValidateBareMetalServerName(input)
		assert.Nil(t, err)
	}
}

func TestValidateBareMetalServerName_InvalidInputs(t *testing.T) {
	invalidInputs := []string{
		"",
		strings.Repeat("a", 64),
		strings.Repeat(strings.Repeat("a", 49)+".", 5) + "bike",
		".leading-period",
		"training-period.",
		"repeated..periods",
		"-leading-hyphen",
	}

	for _, invalidInput := range invalidInputs {
		err := ValidateBareMetalServerName(invalidInput)
		assert.NotNil(t, err)
	}
}

func TestValidationBareMetalUsername_ValidInputs(t *testing.T) {
	validInputs := []string{
		"username",
		"user-name",
		"username$",
		"username123",
	}

	for _, input := range validInputs {
		err := ValidateBareMetalUsername(input)
		assert.Nil(t, err)
	}
}

func TestValidateBareMetalUsername_InvalidInputs(t *testing.T) {
	invalidInputs := []string{
		"",
		"root",
		"ROOT",
		"RooT",
		"user*name",
		"#username",
		"@username",
		"-username",
	}

	for _, input := range invalidInputs {
		err := ValidateBareMetalUsername(input)
		assert.NotNil(t, err)
	}
}

func TestValidateBareMetalPassword_ValidInputs(t *testing.T) {
	validInputs := []string{
		"Testing123",
		"ReallyLongPasswordThatContainsSomeNumbers123",
	}

	for _, input := range validInputs {
		err := ValidateBareMetalPassword(input)
		assert.Nil(t, err)
	}
}

func TestValidateBareMetalPassword_InvalidInput(t *testing.T) {
	invalidInputs := []string{
		"Testing123!",
		"Testing123@",
		"Testing123#",
		"Testing123^",
		"ToShort",
		"@TestBlah",
		"testing123",
		"Testing 123",
		fmt.Sprintf("Testing\n123"),
	}

	for _, input := range invalidInputs {
		err := ValidateBareMetalPassword(input)
		assert.NotNil(t, err)
	}
}
