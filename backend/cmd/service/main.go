package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	datamanager "ss/internal/service/datamanager"
	"ss/internal/service/db"
	"ss/internal/service/handler"
	"time"
)

func main() {
	ctx := context.Background()
	database, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB(ctx)
	dm := datamanager.New(database)
	h := handler.New(dm)

	r := mux.NewRouter()
	r.HandleFunc("/users", h.UsersHandler).Methods("GET")
	r.HandleFunc("/messages/{login}", h.MessagesHandler).Methods("GET")
	r.HandleFunc("/message/send/{login}", h.SendHandler).Methods("POST")
	r.HandleFunc("/auth", h.AuthHandler).Methods("POST")
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
