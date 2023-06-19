package server

import (
	"github.com/nayakunin/gophermart/internal/config"
)

type Server struct {
	Worker          Worker
	Storage         Storage
	TokenService    TokenService
	Cfg             config.Config
	ChecksumService ChecksumService
}

func NewServer(dbStorage Storage, cfg config.Config, w Worker, tokenService TokenService, checksumService ChecksumService) Server {
	go w.Start()

	return Server{
		Worker:          w,
		Storage:         dbStorage,
		Cfg:             cfg,
		TokenService:    tokenService,
		ChecksumService: checksumService,
	}
}
