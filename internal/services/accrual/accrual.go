package accrual

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/nayakunin/gophermart/internal/logger"
)

var (
	ErrNoContent       = errors.New("no content")
	ErrTooManyRequests = errors.New("too many requests")
)

const (
	StatusRegistered = "REGISTERED"
	StatusInvalid    = "INVALID"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
)

type Service struct {
	client *resty.Client
}

func NewAccrualService(url string) *Service {
	client := resty.New()
	client.SetBaseURL(url)

	return &Service{
		client: client,
	}
}

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func (s *Service) GetAccrual(orderID int64) (*resty.Response, *Accrual, error) {
	var accrual Accrual
	resp, err := s.client.R().
		SetResult(&accrual).
		Get(fmt.Sprintf("/api/orders/%d", orderID))

	if err != nil {
		logger.Errorf("failed to get accrual: %v", err)
		return resp, nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusNoContent {
			logger.Errorf("no content")
			return resp, nil, ErrNoContent
		}

		if resp.StatusCode() == http.StatusTooManyRequests {
			logger.Errorf("too many requests")
			return resp, nil, ErrTooManyRequests
		}

		logger.Errorf("unexpected status code: %d", resp.StatusCode())
		return resp, nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return resp, &accrual, nil
}
