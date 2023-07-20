package worker

import (
	stdErrors "errors"

	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/services/accrual"
	"github.com/pkg/errors"
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
	w := &Worker{
		accrual: accrual,
		storage: storage,
		queue:   make(chan Order, MaxRequests),
	}

	go w.Run()

	return w
}

func (w *Worker) AddOrder(userID int64, orderID int64) {
	w.queue <- Order{
		id:     orderID,
		userID: userID,
	}
}

func (w *Worker) Run() {
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
		if stdErrors.Is(err, accrual.ErrNoContent) {
			return errors.Wrap(err, "failed to get accrual")
		}

		return errors.Wrap(err, "failed to get accrual")
	}

	switch accr.Status {
	case accrual.StatusRegistered:
		w.queue <- order
	case accrual.StatusInvalid:
		if err := w.storage.UpdateOrderStatus(order.id, api.OrderStatusINVALID); err != nil {
			return errors.Wrap(err, "failed to update order status")
		}
	case accrual.StatusProcessing:
		w.queue <- order
		if err := w.storage.UpdateOrderStatus(order.id, api.OrderStatusPROCESSING); err != nil {
			return errors.Wrap(err, "failed to update order status")
		}
	case accrual.StatusProcessed:
		if err := w.storage.ProcessOrder(order.userID, order.id, accr.Accrual); err != nil {
			return errors.Wrap(err, "failed to process order")
		}
	}

	return nil
}
