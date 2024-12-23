package main

import (
	"context"
	"log"
	datamanager "ss/internal/service/datamanager"
	"ss/internal/service/db"
	"ss/internal/service/handler"
	"ss/internal/service/router"
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

	r := router.New(h)

	r.Route()
}
