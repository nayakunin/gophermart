package worker

import (
	"errors"
	"time"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/services/accrual"
)

const MaxRequests = 10

type Order struct {
	id     int64
	userID int64
}

type Worker struct {
	accrual AccrualService
	storage Storage
	queue   chan Order
}

func NewWorker(accrual AccrualService, storage Storage) *Worker {
	return &Worker{
		accrual: accrual,
		storage: storage,
		queue:   make(chan Order, MaxRequests),
	}
}

func (w *Worker) AddOrder(userID int64, orderID int64) {
	w.queue <- Order{
		id:     orderID,
		userID: userID,
	}
}

func (w *Worker) Start() {
	for order := range w.queue {
		w.processOrder(order)
	}
}

func (w *Worker) processOrder(order Order) {
	resp, accr, err := w.accrual.GetAccrual(order.id)
	if err != nil {
		if errors.Is(err, accrual.ErrTooManyRequests) {
			retryAfter, err := time.ParseDuration(resp.Header().Get("Retry-After") + "s")
			if err != nil {
				return
			}

			time.Sleep(retryAfter)
			w.queue <- order
		}

		if errors.Is(err, accrual.ErrNoContent) {
			return
		}

		return
	}

	switch accr.Status {
	case accrual.StatusRegistered:
		w.queue <- order
	case accrual.StatusInvalid:
		err := w.storage.UpdateOrderStatus(order.id, api.OrderStatusINVALID)
		if err != nil {
			return
		}
	case accrual.StatusProcessing:
		w.queue <- order
		err := w.storage.UpdateOrderStatus(order.id, api.OrderStatusPROCESSING)
		if err != nil {
			return
		}
	case accrual.StatusProcessed:
		err := w.storage.ProcessOrder(order.userID, order.id, accr.Accrual)
		if err != nil {
			return
		}
	}
}
