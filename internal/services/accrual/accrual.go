package accrual

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nayakunin/gophermart/internal/logger"
)

var (
	ErrNoContent = errors.New("no content")
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

func NewService(url string) *Service {
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
	for {
		var accrual Accrual
		resp, err := s.client.R().
			SetResult(&accrual).
			Get(fmt.Sprintf("/api/orders/%d", orderID))
		if err != nil {
			logger.Errorf("failed to get accrual: %v", err)
			return nil, err
		}

		if resp.StatusCode() != http.StatusOK {
			if resp.StatusCode() == http.StatusNoContent {
				logger.Errorf("no content")
				return nil, ErrNoContent
			}

			if resp.StatusCode() == http.StatusTooManyRequests {
				logger.Errorf("too many requests")
				retryAfter, err := time.ParseDuration(resp.Header().Get("Retry-After") + "s")
				if err != nil {
					return nil, err
				}
				time.Sleep(retryAfter)
				continue
			}

			logger.Errorf("unexpected status code: %d", resp.StatusCode())
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		return &accrual, nil
	}
}
