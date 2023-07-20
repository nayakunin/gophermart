package server

import (
	"github.com/nayakunin/gophermart/internal/config"
	"github.com/nayakunin/gophermart/internal/storage"
)

func NewMockServer(worker Worker, storage Storage, cfg config.Config, tokenService TokenService) Server {
	return Server{
		Worker:       worker,
		Storage:      storage,
		Cfg:          cfg,
		TokenService: tokenService,
	}
}

type MockStorage struct {
	CreateUserResponse     int64
	CreateUserError        error
	GetUserIDResponse      int64
	GetUserIDError         error
	SaveOrderError         error
	GetOrdersResponse      []storage.Order
	GetOrdersError         error
	GetBalanceResponse     storage.Balance
	GetBalanceError        error
	WithdrawError          error
	GetWithdrawalsResponse []storage.Transaction
	GetWithdrawalsError    error
}

func (m MockStorage) CreateUser(_ string, _ string) (int64, error) {
	return m.CreateUserResponse, m.CreateUserError
}

func (m MockStorage) GetUserID(_ string, _ string) (int64, error) {
	return m.GetUserIDResponse, m.GetUserIDError
}

func (m MockStorage) SaveOrder(_ int64, _ int64) error {
	return m.SaveOrderError
}

func (m MockStorage) GetOrders(_ int64) ([]storage.Order, error) {
	return m.GetOrdersResponse, m.GetOrdersError
}

func (m MockStorage) GetBalance(_ int64) (storage.Balance, error) {
	return m.GetBalanceResponse, m.GetBalanceError
}

func (m MockStorage) Withdraw(_ int64, _ int64, _ float32) error {
	return m.WithdrawError
}

func (m MockStorage) GetWithdrawals(_ int64) ([]storage.Transaction, error) {
	return m.GetWithdrawalsResponse, m.GetWithdrawalsError
}

type StorageParams struct {
	CreateUserResponse     int64
	CreateUserError        error
	GetUserIDResponse      int64
	GetUserIDError         error
	SaveOrderError         error
	GetOrdersResponse      []storage.Order
	GetOrdersError         error
	GetBalanceResponse     storage.Balance
	GetBalanceError        error
	WithdrawError          error
	GetWithdrawalsResponse []storage.Transaction
	GetWithdrawalsError    error
}

func NewMockStorage(params StorageParams) Storage {
	return &MockStorage{
		CreateUserResponse:     params.CreateUserResponse,
		CreateUserError:        params.CreateUserError,
		GetUserIDResponse:      params.GetUserIDResponse,
		GetUserIDError:         params.GetUserIDError,
		SaveOrderError:         params.SaveOrderError,
		GetOrdersResponse:      params.GetOrdersResponse,
		GetOrdersError:         params.GetOrdersError,
		GetBalanceResponse:     params.GetBalanceResponse,
		GetBalanceError:        params.GetBalanceError,
		WithdrawError:          params.WithdrawError,
		GetWithdrawalsResponse: params.GetWithdrawalsResponse,
		GetWithdrawalsError:    params.GetWithdrawalsError,
	}
}

type MockTokenService struct {
	CreateTokenResponse string
	CreateTokenError    error
}

func (m MockTokenService) CreateToken(_ config.Config, _ int64) (string, error) {
	return m.CreateTokenResponse, m.CreateTokenError
}

type TokenServiceParams struct {
	CreateTokenResponse string
	CreateTokenError    error
}

func NewMockTokenService(params TokenServiceParams) TokenService {
	return &MockTokenService{
		CreateTokenResponse: params.CreateTokenResponse,
		CreateTokenError:    params.CreateTokenError,
	}
}
