package main

import (
	"ssd-coursework/routes/callback"

	"github.com/gorilla/mux"
)

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/callback", callback.CallbackHandler)
}
