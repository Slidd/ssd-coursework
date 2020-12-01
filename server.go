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
	// Redirect HTTP to HTTPS
	go http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://localhost:3000/"+r.URL.String(), http.StatusMovedPermanently)
	}))
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
		negroni.HandlerFunc(middlewares.AuthorizedToAccess),
		negroni.Wrap(http.HandlerFunc(crud.Show)),
	))
	r.Handle("/new", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.HandlerFunc(middlewares.AuthorizedToAccess),
		negroni.Wrap(http.HandlerFunc(crud.New)),
	))
	r.Handle("/newTicket", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.HandlerFunc(middlewares.AuthorizedToAccess),
		negroni.Wrap(http.HandlerFunc(crud.NewTicket)),
	))
	r.Handle("/edit", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.HandlerFunc(middlewares.AuthorizedToAccess),
		negroni.Wrap(http.HandlerFunc(crud.Edit)),
	))
	r.Handle("/addComment", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.HandlerFunc(middlewares.AuthorizedToAccess),
		negroni.Wrap(http.HandlerFunc(crud.AddComment)),
	))
	r.Handle("/update", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.HandlerFunc(middlewares.AuthorizedToAccess),
		negroni.Wrap(http.HandlerFunc(crud.UpdateTicket)),
	))
	// r.Handle("/delete", negroni.New(
	// 	negroni.HandlerFunc(middlewares.IsAuthenticated),
	// 	negroni.Wrap(http.HandlerFunc(crud.Delete)),
	// ))

	r.HandleFunc("/login", login.LoginHandler)
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

	// Files we want to serve to the web application (i.e. files that will be used during execution)
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.Handle("/", r)
	log.Print("Server listening on https://localhost:3000/")

	log.Fatal(http.ListenAndServeTLS("127.0.0.1:3000", "../localhost.crt", "../localhost.key", nil))
}
