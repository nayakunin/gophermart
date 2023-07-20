package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/nayakunin/gophermart/internal/config"
	api "github.com/nayakunin/gophermart/internal/generated"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
)

func TestGetAPIUserOrders(t *testing.T) {
	userID := int64(1)
	type args struct {
		orders []storage.Order
		err    error
	}

	accrual := float32(1)
	uploadedAt := time.Time{}

	tests := []struct {
		name string
		args args
		want *api.Response
	}{
		{
			name: "OK",
			args: args{
				orders: []storage.Order{{
					ID:         1,
					Status:     api.OrderStatusNEW,
					UploadedAt: uploadedAt,
					Accrual:    &accrual,
				}},
			},
			want: api.GetAPIUserOrdersJSON200Response([]api.GetOrdersOrder{{
				Number:     "1",
				Status:     api.OrderStatusNEW,
				UploadedAt: uploadedAt,
				Accrual:    &accrual,
			}}),
		},
		{
			name: "Error",
			args: args{
				err: fmt.Errorf("error"),
			},
			want: &api.Response{Code: http.StatusInternalServerError},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				AuthKey: "auth_key",
			}

			storageParams := StorageParams{}

			if tt.args.orders != nil {
				storageParams.GetOrdersResponse = tt.args.orders
			}

			if tt.args.err != nil {
				storageParams.GetOrdersError = tt.args.err
			}

			mockStorage := NewMockStorage(storageParams)
			server := NewMockServer(nil, mockStorage, cfg, nil)
			logger.Mock()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
			r = r.WithContext(context.WithValue(r.Context(), cfg.AuthKey, userID))

			if got := server.GetAPIUserOrders(w, r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAPIUserOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}
