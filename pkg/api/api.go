package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"gitlab.com/slirx/newproj/pkg/tracer"
)

// MessageType represents message type.
type MessageType string

const (
	MessageTypeSuccess MessageType = "success"
	MessageTypeError   MessageType = "error"
)

var InternalError = errors.New("oops, something went wrong")
var RequestError = errors.New("invalid request data")

type requestError struct {
	Err error
}

func (r requestError) Error() string {
	return r.Err.Error()
}

type accessError struct {
	Err error
}

func (r accessError) Error() string {
	return r.Err.Error()
}

type notFoundError struct {
	Err error
}

func (r notFoundError) Error() string {
	return r.Err.Error()
}

func NewRequestError(err error) error {
	return requestError{Err: err}
}

func NewAccessError(err error) error {
	return accessError{Err: err}
}

func NewNotFoundError(err error) error {
	return notFoundError{Err: err}
}

// Response represents fields which should be in every response.
type Response struct {
	RequestID string `json:"request_id"`
}

// MessageResponse represents fields which should be in responses with text messages.
type MessageResponse struct {
	RequestID string      `json:"request_id"`
	Type      MessageType `json:"type,omitempty"`
	Message   string      `json:"message,omitempty"`
}

// DataResponse represents fields which should be in responses with data objects.
type DataResponse struct {
	RequestID string      `json:"request_id"`
	Type      MessageType `json:"type,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// todo move these functions to interface, so they can be mocked

// NewResponse creates new instance of Response.
func NewResponse(requestID string) Response {
	return Response{RequestID: requestID}
}

// NewMessageResponse creates new instance of MessageResponse.
func NewMessageResponse(requestID string, messageType MessageType, message string) MessageResponse {
	return MessageResponse{
		RequestID: requestID,
		Type:      messageType,
		Message:   message,
	}
}

// NewDataResponse creates new instance of DataResponse.
func NewDataResponse(requestID string, messageType MessageType, data interface{}) DataResponse {
	return DataResponse{
		RequestID: requestID,
		Type:      messageType,
		Data:      data,
	}
}

func GetErrorResponseFields(err error) (int, string) {
	if errors.Is(err, RequestError) {
		return http.StatusBadRequest, err.Error()
	}

	switch t := err.(type) {
	case requestError:
		return http.StatusBadRequest, t.Error()
	case accessError:
		return http.StatusForbidden, t.Error()
	case notFoundError:
		return http.StatusNotFound, t.Error()
	default:
		return http.StatusInternalServerError, InternalError.Error()
	}
}

type ResponseBuilder interface {
	ErrorResponse(ctx context.Context, w http.ResponseWriter, err error)
	MessageResponse(ctx context.Context, w http.ResponseWriter, message string)
	DataResponse(ctx context.Context, w http.ResponseWriter, data interface{})
}

type responseBuilder struct {
	Tracer tracer.Tracer
}

func (r responseBuilder) ErrorResponse(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	responseCode, responseMessage := GetErrorResponseFields(err)
	w.WriteHeader(responseCode)

	response := NewMessageResponse(r.Tracer.RequestID(ctx), MessageTypeError, responseMessage)
	_ = json.NewEncoder(w).Encode(response)
}

func (r responseBuilder) MessageResponse(ctx context.Context, w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")

	response := NewMessageResponse(r.Tracer.RequestID(ctx), MessageTypeSuccess, message)
	_ = json.NewEncoder(w).Encode(response)
}

func (r responseBuilder) DataResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	response := NewDataResponse(r.Tracer.RequestID(ctx), MessageTypeSuccess, data)
	_ = json.NewEncoder(w).Encode(response)
}

func NewResponseBuilder(t tracer.Tracer) ResponseBuilder {
	return responseBuilder{Tracer: t}
}
