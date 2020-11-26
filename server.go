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

	r.Handle("/", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.HandlerFunc(middlewares.AuthorizedToAccess),
		negroni.Wrap(http.HandlerFunc(crud.Index)),
	))
	r.Handle("/logout", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(logout.LogoutHandler)),
	))
	r.Handle("/show", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(crud.Show)),
	))
	r.Handle("/new", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(crud.New)),
	))
	r.Handle("/edit", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(crud.Edit)),
	))
	r.Handle("/addComment", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(crud.AddComment)),
	))
	r.Handle("/update", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(crud.Update)),
	))
	r.Handle("/delete", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(crud.Delete)),
	))

	r.HandleFunc("/login", login.LoginHandler)
	r.HandleFunc("/logout", logout.LogoutHandler)
	r.HandleFunc("/callback", callback.CallbackHandler)
	r.Handle("/user", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(user.UserHandler)),
	))

	// http.HandleFunc("/index", crud.Index)
	// http.HandleFunc("/show", crud.Show)
	// http.HandleFunc("/new", crud.New)
	// http.HandleFunc("/edit", crud.Edit)
	// http.HandleFunc("/insert", crud.Insert)
	// http.HandleFunc("/update", crud.Update)
	// http.HandleFunc("/delete", crud.Delete)

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.Handle("/", r)
	log.Print("Server listening on http://localhost:3000/")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
