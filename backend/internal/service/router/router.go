package router

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"ss/internal/service/handler"
)

type Router interface {
	Route()
}

type router struct {
	handler handler.Handler
}

func New(handler handler.Handler) Router {
	return &router{handler: handler}
}

func (rt *router) Route() {
	r := mux.NewRouter()
	r.HandleFunc("/users", rt.handler.UsersHandler).Methods("GET")
	r.HandleFunc("/messages/{login}", rt.handler.MessagesHandler).Methods("GET")
	r.HandleFunc("/message/send/{login}", rt.handler.SendHandler).Methods("POST")
	r.HandleFunc("/auth", rt.handler.AuthHandler).Methods("POST")
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         serveAddr,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	log.Fatal(srv.ListenAndServe())
}
