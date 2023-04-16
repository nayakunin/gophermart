package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nayakunin/gophermart/internal/config"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/middlewares"
	"github.com/nayakunin/gophermart/internal/server"
	"github.com/nayakunin/gophermart/internal/storage"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbStorage, err := storage.NewDBStorage(c.DataBaseURI)
	if err != nil {
		log.Fatal(err)
	}

	apiImpl := server.NewServer(dbStorage, *c)

	r.Mount("/", api.Handler(apiImpl, api.WithMiddleware("auth", middlewares.Auth(*c))))

	log.Fatal(http.ListenAndServe(c.RunAddress, r))
}
