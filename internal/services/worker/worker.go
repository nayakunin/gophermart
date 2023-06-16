package worker

import (
	"errors"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
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
		go func(order Order) {
			if err := w.processOrder(order); err != nil {
				logger.Errorf("failed to process order: %v", err)
			}
		}(order)
	}
}

func (w *Worker) processOrder(order Order) error {
	accr, err := w.accrual.GetAccrual(order.id)
	if err != nil {
		if errors.Is(err, accrual.ErrNoContent) {
			return err
		}

		return err
	}

	switch accr.Status {
	case accrual.StatusRegistered:
		w.queue <- order
	case accrual.StatusInvalid:
		err := w.storage.UpdateOrderStatus(order.id, api.OrderStatusINVALID)
		if err != nil {
			return err
		}
	case accrual.StatusProcessing:
		w.queue <- order
		err := w.storage.UpdateOrderStatus(order.id, api.OrderStatusPROCESSING)
		if err != nil {
			return err
		}
	case accrual.StatusProcessed:
		err := w.storage.ProcessOrder(order.userID, order.id, accr.Accrual)
		if err != nil {
			return err
		}
	}

	return nil
}
