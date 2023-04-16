// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.2.2 DO NOT EDIT.
package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Defines values for OrderStatus.
var (
	UnknownOrderStatus = OrderStatus{}

	OrderStatusINVALID = OrderStatus{"INVALID"}

	OrderStatusNEW = OrderStatus{"NEW"}

	OrderStatusPROCESSED = OrderStatus{"PROCESSED"}

	OrderStatusPROCESSING = OrderStatus{"PROCESSING"}
)

// Balance defines model for Balance.
type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

// BalanceWithdrawRequest defines model for BalanceWithdrawRequest.
type BalanceWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

// GetOrdersOrder defines model for GetOrdersOrder.
type GetOrdersOrder struct {
	Accrual   *float32    `json:"accrual,omitempty"`
	Number    string      `json:"number"`
	Status    OrderStatus `json:"status"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// GetOrdersReply defines model for GetOrdersReply.
type GetOrdersReply struct {
	Orders []GetOrdersOrder `json:"orders"`
}

// GetUserWithdrawalsReply defines model for GetUserWithdrawalsReply.
type GetUserWithdrawalsReply []GetUserWithdrawalsReplyItem

// GetUserWithdrawalsReplyItem defines model for GetUserWithdrawalsReplyItem.
type GetUserWithdrawalsReplyItem struct {
	Order       string    `json:"order"`
	ProcessedAt time.Time `json:"processed_at"`
	Sum         float32   `json:"sum"`
}

// RegisterUserRequest defines model for RegisterUserRequest.
type RegisterUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// OrderStatus defines model for OrderStatus.
type OrderStatus struct {
	value string
}

func (t *OrderStatus) ToValue() string {
	return t.value
}
func (t *OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}
func (t *OrderStatus) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return t.FromValue(value)
}
func (t *OrderStatus) FromValue(value string) error {
	switch value {

	case OrderStatusINVALID.value:
		t.value = value
		return nil

	case OrderStatusNEW.value:
		t.value = value
		return nil

	case OrderStatusPROCESSED.value:
		t.value = value
		return nil

	case OrderStatusPROCESSING.value:
		t.value = value
		return nil

	}
	return fmt.Errorf("unknown enum value: %v", value)
}

// PostAPIUserBalanceWithdrawJSONBody defines parameters for PostAPIUserBalanceWithdraw.
type PostAPIUserBalanceWithdrawJSONBody BalanceWithdrawRequest

// PostAPIUserLoginJSONBody defines parameters for PostAPIUserLogin.
type PostAPIUserLoginJSONBody RegisterUserRequest

// PostAPIUserRegisterJSONBody defines parameters for PostAPIUserRegister.
type PostAPIUserRegisterJSONBody RegisterUserRequest

// PostAPIUserBalanceWithdrawJSONRequestBody defines body for PostAPIUserBalanceWithdraw for application/json ContentType.
type PostAPIUserBalanceWithdrawJSONRequestBody PostAPIUserBalanceWithdrawJSONBody

// Bind implements render.Binder.
func (PostAPIUserBalanceWithdrawJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostAPIUserLoginJSONRequestBody defines body for PostAPIUserLogin for application/json ContentType.
type PostAPIUserLoginJSONRequestBody PostAPIUserLoginJSONBody

// Bind implements render.Binder.
func (PostAPIUserLoginJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostAPIUserRegisterJSONRequestBody defines body for PostAPIUserRegister for application/json ContentType.
type PostAPIUserRegisterJSONRequestBody PostAPIUserRegisterJSONBody

// Bind implements render.Binder.
func (PostAPIUserRegisterJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// GetAPIUserBalanceJSON200Response is a constructor method for a GetAPIUserBalance response.
// A *Response is returned with the configured status code and content type from the spec.
func GetAPIUserBalanceJSON200Response(body Balance) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetAPIUserOrdersJSON200Response is a constructor method for a GetAPIUserOrders response.
// A *Response is returned with the configured status code and content type from the spec.
func GetAPIUserOrdersJSON200Response(body GetOrdersReply) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetAPIUserWithdrawalsJSON200Response is a constructor method for a GetAPIUserWithdrawals response.
// A *Response is returned with the configured status code and content type from the spec.
func GetAPIUserWithdrawalsJSON200Response(body GetUserWithdrawalsReply) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get user balance
	// (GET /api/user/balance)
	GetAPIUserBalance(w http.ResponseWriter, r *http.Request) *Response
	// Withdraw user balance
	// (POST /api/user/balance/withdraw)
	PostAPIUserBalanceWithdraw(w http.ResponseWriter, r *http.Request) *Response
	// Login user
	// (POST /api/user/login)
	PostAPIUserLogin(w http.ResponseWriter, r *http.Request) *Response
	// Get user orders
	// (GET /api/user/orders)
	GetAPIUserOrders(w http.ResponseWriter, r *http.Request) *Response
	// Create user order
	// (POST /api/user/orders)
	PostAPIUserOrders(w http.ResponseWriter, r *http.Request) *Response
	// Register user
	// (POST /api/user/register)
	PostAPIUserRegister(w http.ResponseWriter, r *http.Request) *Response
	// Get user withdrawals
	// (GET /api/user/withdrawals)
	GetAPIUserWithdrawals(w http.ResponseWriter, r *http.Request) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// GetAPIUserBalance operation middleware
func (siw *ServerInterfaceWrapper) GetAPIUserBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetAPIUserBalance(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	// Operation specific middleware
	handler = siw.Middlewares["auth"](handler).ServeHTTP

	handler(w, r.WithContext(ctx))
}

// PostAPIUserBalanceWithdraw operation middleware
func (siw *ServerInterfaceWrapper) PostAPIUserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostAPIUserBalanceWithdraw(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	// Operation specific middleware
	handler = siw.Middlewares["auth"](handler).ServeHTTP

	handler(w, r.WithContext(ctx))
}

// PostAPIUserLogin operation middleware
func (siw *ServerInterfaceWrapper) PostAPIUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostAPIUserLogin(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetAPIUserOrders operation middleware
func (siw *ServerInterfaceWrapper) GetAPIUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetAPIUserOrders(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	// Operation specific middleware
	handler = siw.Middlewares["auth"](handler).ServeHTTP

	handler(w, r.WithContext(ctx))
}

// PostAPIUserOrders operation middleware
func (siw *ServerInterfaceWrapper) PostAPIUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostAPIUserOrders(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	// Operation specific middleware
	handler = siw.Middlewares["auth"](handler).ServeHTTP

	handler(w, r.WithContext(ctx))
}

// PostAPIUserRegister operation middleware
func (siw *ServerInterfaceWrapper) PostAPIUserRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostAPIUserRegister(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetAPIUserWithdrawals operation middleware
func (siw *ServerInterfaceWrapper) GetAPIUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetAPIUserWithdrawals(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	// Operation specific middleware
	handler = siw.Middlewares["auth"](handler).ServeHTTP

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      map[string]func(http.Handler) http.Handler
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:     "/",
		BaseRouter:  chi.NewRouter(),
		Middlewares: make(map[string]func(http.Handler) http.Handler),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		Middlewares:      options.Middlewares,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	middlewares := []string{"auth"}
	for _, m := range middlewares {
		if _, ok := wrapper.Middlewares[m]; !ok {
			panic("goapi-gen: could not find tagged middleware " + m)
		}
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/api/user/balance", wrapper.GetAPIUserBalance)
		r.Post("/api/user/balance/withdraw", wrapper.PostAPIUserBalanceWithdraw)
		r.Post("/api/user/login", wrapper.PostAPIUserLogin)
		r.Get("/api/user/orders", wrapper.GetAPIUserOrders)
		r.Post("/api/user/orders", wrapper.PostAPIUserOrders)
		r.Post("/api/user/register", wrapper.PostAPIUserRegister)
		r.Get("/api/user/withdrawals", wrapper.GetAPIUserWithdrawals)

	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithMiddleware(key string, middleware func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares[key] = middleware
	}
}

func WithMiddlewares(middlewares map[string]func(http.Handler) http.Handler) ServerOption {
	return func(s *ServerOptions) {
		s.Middlewares = middlewares
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}