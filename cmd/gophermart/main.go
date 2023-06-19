package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nayakunin/gophermart/internal/config"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/middlewares"
	"github.com/nayakunin/gophermart/internal/server"
	"github.com/nayakunin/gophermart/internal/services/accrual"
	"github.com/nayakunin/gophermart/internal/services/checksum"
	"github.com/nayakunin/gophermart/internal/services/token"
	"github.com/nayakunin/gophermart/internal/services/worker"
	"github.com/nayakunin/gophermart/internal/storage"
)

func main() {
	logger.Init()

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

	accrualService := accrual.NewService(c.AccrualSystemAddress)
	checkSumService := checksum.NewService()
	tokenService := token.NewService()
	w := worker.NewWorker(accrualService, dbStorage)

	apiImpl := server.NewServer(dbStorage, *c, w, tokenService, checkSumService)

	/**
	 * Middlewares are only applied to the routes that are specified in the api/schema.yaml
	 */
	r.Mount("/", api.Handler(apiImpl, api.WithMiddleware("token", middlewares.Auth(*c))))

	log.Fatal(http.ListenAndServe(c.RunAddress, r))
}
