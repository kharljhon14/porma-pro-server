package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kharljhon14/porma-pro-server/cmd/api"
)

func main() {
	_, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal("cannot connect to DB: ", err)
	}

	server, err := api.NewServer()
	if err != nil {
		log.Fatal("cannot create new server: ", err)
	}

	err = server.Start(os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
