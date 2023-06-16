package server

import (
	"github.com/nayakunin/gophermart/internal/config"
)

type Server struct {
	Worker  Worker
	Storage Storage
	Cfg     config.Config
}

func NewServer(dbStorage Storage, cfg config.Config, w Worker) Server {
	go w.Start()

	return Server{
		Worker:  w,
		Storage: dbStorage,
		Cfg:     cfg,
	}
}
