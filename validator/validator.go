package validator

import (
	"net/http"
	"ssd-coursework/routes/user"
)

// CanChangeTicketStatus Checks if the user has the correct role to change the ticket status
func CanChangeTicketStatus(w http.ResponseWriter, r *http.Request, previousStatus string, currentStatus string) bool {
	if previousStatus != currentStatus {
		if user.IsDeveloper(w, r) || user.IsTester(w, r) {
			return true
		}
		return false
	}
	return true
}
