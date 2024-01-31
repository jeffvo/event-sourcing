package main

import (
	"net/http"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/jeffvo/event-sourcing/api/routes"
)

func main() {

	db := InitDatabase()
	SetupServer(db)
}

func InitDatabase() *esdb.Client {

	settings, err := esdb.ParseConnectionString("esdb+discover://localhost:2113?tls=false")

	if err != nil {
		panic(err)
	}

	db, err := esdb.NewClient(settings)

	if err != nil {
		panic(err)
	}

	return db
}

func SetupServer(db *esdb.Client) {
	var router = routes.RegisterHTTPEndpoints(db)

	var server = &http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	server.ListenAndServe()

	db.Close()
}
