package middlewares

import (
	"log"
	"net/http"

	"ssd-coursework/app"
	"ssd-coursework/routes/user"
)

// IsAuthenticated Check if the user is authenticated
func IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		next(w, r)
	}
}

func AuthorizedToAccess(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if user.IsDeveloper(w, r) || user.IsClient(w, r) || user.IsTester(w, r) {
		next(w, r)
	} else {
		http.Error(w, "Not Authorized to use application", http.StatusInternalServerError)
		return
	}
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		log.Fatal("failed to delete session", err)
	}
}
