package accrual

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
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

func (s *Service) GetAccrual(orderID int64) (*Accrual, error) {
	var accrual Accrual
	req, err := s.client.R().
		SetResult(&accrual).
		Get(fmt.Sprintf("/api/orders/%d", orderID))

	if err != nil {
		return nil, err
	}

	if req.StatusCode() != http.StatusOK {
		if req.StatusCode() == http.StatusNoContent {
			return nil, ErrNoContent
		}

		if req.StatusCode() == http.StatusTooManyRequests {
			return nil, ErrTooManyRequests
		}

		return nil, fmt.Errorf("unexpected status code: %d", req.StatusCode())
	}

	return &accrual, nil
}
