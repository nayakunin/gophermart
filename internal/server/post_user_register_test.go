package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nayakunin/gophermart/internal/config"
	"github.com/nayakunin/gophermart/internal/logger"
	"github.com/nayakunin/gophermart/internal/storage"
)

func TestPostAPIUserRegister(t *testing.T) {
	userID := int64(1)
	type args struct {
		body            string
		createUserError error
		tokenResponse   string
		tokenError      error
	}

	type want struct {
		code   int
		cookie string
	}

	cookie := "token"

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "OK",
			args: args{
				body:          `{"login": "1", "password": "1"}`,
				tokenResponse: "token",
			},
			want: want{
				code:   http.StatusOK,
				cookie: fmt.Sprintf("Authentication=%s", cookie),
			},
		},
		{
			name: "Unable to read body",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "Unable to unmarshal body",
			args: args{
				body: `{"login": "1", "password": "1"`,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "No login",
			args: args{
				body: `{"login": "", "password": "1"}`,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "No password",
			args: args{
				body: `{"login": "1", "password": ""}`,
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "User already exists",
			args: args{
				body:            `{"login": "1", "password": "1"}`,
				createUserError: storage.ErrUserExists,
			},
			want: want{
				code: http.StatusConflict,
			},
		},
		{
			name: "Unknown storage error",
			args: args{
				body:            `{"login": "1", "password": "1"}`,
				createUserError: fmt.Errorf("error"),
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "Auth token generation error",
			args: args{
				body:       `{"login": "1", "password": "1"}`,
				tokenError: fmt.Errorf("error"),
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				AuthKey: "auth_key",
			}

			storageParams := StorageParams{}

			if tt.args.createUserError != nil {
				storageParams.CreateUserError = tt.args.createUserError
			}

			tokenServiceParams := TokenServiceParams{}

			if tt.args.tokenResponse != "" {
				tokenServiceParams.CreateTokenResponse = tt.args.tokenResponse
			}

			if tt.args.tokenError != nil {
				tokenServiceParams.CreateTokenError = tt.args.tokenError
			}

			mockStorage := NewMockStorage(storageParams)
			mockTokenService := NewMockTokenService(tokenServiceParams)
			server := NewMockServer(nil, mockStorage, cfg, mockTokenService)
			logger.Mock()

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/api/user/register", strings.NewReader(tt.args.body))
			r = r.WithContext(context.WithValue(r.Context(), cfg.AuthKey, userID))

			got := server.PostAPIUserRegister(w, r)

			if got.Code != tt.want.code {
				t.Errorf("PostAPIUserRegister() = %v, want %v", got, tt.want)
			}

			if got.Code == http.StatusOK {
				cookie := w.Header().Get("Set-Cookie")
				if cookie != tt.want.cookie {
					t.Errorf("PostAPIUserRegister() cookie = %v, want %v", cookie, tt.want.cookie)
				}
			}
		})
	}
}
