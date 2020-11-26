package main

import (
	"log"
	"net/http"
	"ssd-coursework/routes/callback"
	"ssd-coursework/routes/crud"
	"ssd-coursework/routes/login"
	"ssd-coursework/routes/logout"
	"ssd-coursework/routes/middlewares"
	"ssd-coursework/routes/user"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// StartServer start server and routes
func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", crud.Index)
	r.HandleFunc("/login", login.LoginHandler)
	r.HandleFunc("/logout", logout.LogoutHandler)
	r.HandleFunc("/callback", callback.CallbackHandler)
	r.Handle("/user", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(user.UserHandler)),
	))

	http.HandleFunc("/index", crud.Index)
	http.HandleFunc("/show", crud.Show)
	http.HandleFunc("/new", crud.New)
	http.HandleFunc("/edit", crud.Edit)
	http.HandleFunc("/insert", crud.Insert)
	http.HandleFunc("/update", crud.Update)
	http.HandleFunc("/delete", crud.Delete)

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.Handle("/", r)
	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
