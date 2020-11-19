package home

import (
	"net/http"
	"ssd-coursework/routes/templates"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.RenderTemplate(w, "home", nil)
}
