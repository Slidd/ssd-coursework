package validator

import (
	"net/http"
	"reflect"
	"ssd-coursework/routes/user"
)

// CanChangeTicketStatus Checks if the user has the correct role to change the ticket status
func CanChangeTicketStatus(w http.ResponseWriter, r *http.Request, previousStatus string, currentStatus string) bool {
	role := user.GetUserRole(w, r)
	if previousStatus != currentStatus {
		if isDeveloper(role) || isTester(role) {
			return true
		}
		return false
	}
	return true
}

func CanAccessApplication(user interface{}) bool {
	if isDeveloper(user) || isClient(user) || isTester(user) {
		return true
	}
	return false
}

// IsDeveloper checks if the logged in user has the Developer role
func isDeveloper(user interface{}) bool {
	if userHasRole(user, "Developer") {
		return true
	}
	return false
}

// IsClient checks if the logged in user has the Client role
func isClient(user interface{}) bool {

	if userHasRole(user, "Client") {
		return true
	}
	return false
}

// IsTester check if the logged in user has the Tester role
func isTester(user interface{}) bool {

	if userHasRole(user, "Tester") {
		return true
	}
	return false
}

// Helper function to check the values inside the user interface
// Interface is a slice, but first need to reflect this type so Go
// knows how to handle it.
func userHasRole(user interface{}, role string) bool {
	switch reflect.TypeOf(user).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(user)

		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == role {
				return true
			}
		}
	}
	return false
}
