package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nayakunin/gophermart/internal/config"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/middlewares"
	"github.com/nayakunin/gophermart/internal/server"
	"github.com/nayakunin/gophermart/internal/storage"
)

func main() {
	r := chi.NewRouter()

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbStorage, err := storage.NewDBStorage(c.DataBaseURI)
	if err != nil {
		log.Fatal(err)
	}

	apiImpl := server.NewServer(dbStorage)

	r.Mount("/", api.Handler(apiImpl, api.WithMiddleware("auth", middlewares.Auth)))

	log.Fatal(http.ListenAndServe(c.ServerAddress, r))
}
