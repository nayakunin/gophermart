// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.2.2 DO NOT EDIT.
package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

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

// Order defines model for Order.
type Order struct {
	Accrual *int        `json:"accrual,omitempty"`
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
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

// GetAPIOrdersNumberJSON200Response is a constructor method for a GetAPIOrdersNumber response.
// A *Response is returned with the configured status code and content type from the spec.
func GetAPIOrdersNumberJSON200Response(body Order) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}