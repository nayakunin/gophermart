package server

import (
	"github.com/nayakunin/gophermart/internal/config"
	"github.com/nayakunin/gophermart/internal/services/accrual"
	"github.com/nayakunin/gophermart/internal/storage"
)

type Server struct {
	Worker  *Worker
	Accrual *accrual.Service
	Storage *storage.DBStorage
	Cfg     config.Config
}

func NewServer(dbStorage *storage.DBStorage, cfg config.Config, accrualService *accrual.Service) Server {
	worker := NewWorker(accrualService, dbStorage)

	go worker.Start()

	return Server{
		Worker:  worker,
		Accrual: accrualService,
		Storage: dbStorage,
		Cfg:     cfg,
	}
}
