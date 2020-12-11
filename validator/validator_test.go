package validator

import (
	"testing"
)

func TestIsDeveloperHasRole(t *testing.T) {
	var userRole []interface{}
	userRole = append(userRole, "Developer")

	expectedOutput := true

	actualOutput := isDeveloper(userRole)
	if actualOutput != expectedOutput {
		t.Errorf("Actual Output was incorrect, expected %v but got %v", expectedOutput, actualOutput)
	}
}

func TestIsDeveloperNotRole(t *testing.T) {
	var userRole []interface{}
	userRole = append(userRole, "dummy")

	expectedOutput := false

	actualOutput := isDeveloper(userRole)
	if actualOutput != expectedOutput {
		t.Errorf("Actual Output was incorrect, expected %v but got %v", expectedOutput, actualOutput)
	}
}

func TestIsTesterHasRole(t *testing.T) {
	var userRole []interface{}
	userRole = append(userRole, "Tester")

	expectedOutput := true

	actualOutput := isTester(userRole)
	if actualOutput != expectedOutput {
		t.Errorf("Actual Output was incorrect, expected %v but got %v", expectedOutput, actualOutput)
	}
}

func TestIsTesterNotRole(t *testing.T) {
	var userRole []interface{}
	userRole = append(userRole, "dummy")

	expectedOutput := false

	actualOutput := isTester(userRole)
	if actualOutput != expectedOutput {
		t.Errorf("Actual Output was incorrect, expected %v but got %v", expectedOutput, actualOutput)
	}
}
func TestIsClientHasRole(t *testing.T) {
	var userRole []interface{}
	userRole = append(userRole, "Client")

	expectedOutput := true

	actualOutput := isClient(userRole)
	if actualOutput != expectedOutput {
		t.Errorf("Actual Output was incorrect, expected %v but got %v", expectedOutput, actualOutput)
	}
}

func TestIsClientNotRole(t *testing.T) {
	var userRole []interface{}
	userRole = append(userRole, "dummy")

	expectedOutput := false

	actualOutput := isClient(userRole)
	if actualOutput != expectedOutput {
		t.Errorf("Actual Output was incorrect, expected %v but got %v", expectedOutput, actualOutput)
	}
}
