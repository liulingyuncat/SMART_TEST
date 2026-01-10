package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// HandlerFunc is a function that handles a JSON-RPC method call.
type HandlerFunc func(params json.RawMessage) (interface{}, error)

// HandlerFuncWithContext is a function that handles a JSON-RPC method call with context.
type HandlerFuncWithContext func(ctx context.Context, params json.RawMessage) (interface{}, error)

// MessageRouter routes JSON-RPC requests to their handlers.
type MessageRouter struct {
	handlers    map[string]HandlerFunc
	handlersCtx map[string]HandlerFuncWithContext
	mu          sync.RWMutex
}

// NewMessageRouter creates a new MessageRouter.
func NewMessageRouter() *MessageRouter {
	return &MessageRouter{
		handlers:    make(map[string]HandlerFunc),
		handlersCtx: make(map[string]HandlerFuncWithContext),
	}
}

// Register registers a handler for a method.
func (r *MessageRouter) Register(method string, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[method] = handler
}

// RegisterWithContext registers a handler with context for a method.
func (r *MessageRouter) RegisterWithContext(method string, handler HandlerFuncWithContext) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlersCtx[method] = handler
}

// Route routes a request to the appropriate handler and returns a response.
func (r *MessageRouter) Route(req *Request) *Response {
	r.mu.RLock()
	handler, exists := r.handlers[req.Method]
	handlerCtx, existsCtx := r.handlersCtx[req.Method]
	r.mu.RUnlock()

	if !exists && !existsCtx {
		return BuildErrorResponse(
			req.ID,
			MethodNotFound,
			fmt.Sprintf("Method not found: %s", req.Method),
			nil,
		)
	}

	var result interface{}
	var err error

	// Prefer context-aware handler
	if existsCtx {
		ctx := req.Context
		if ctx == nil {
			ctx = context.Background()
		}
		result, err = handlerCtx(ctx, req.Params)
	} else {
		result, err = handler(req.Params)
	}

	if err != nil {
		// Check if the error is already a protocol.Error
		if protoErr, ok := err.(*Error); ok {
			return BuildErrorResponse(req.ID, protoErr.Code, protoErr.Message, protoErr.Data)
		}
		// Otherwise, wrap as internal error
		return BuildErrorResponse(req.ID, InternalError, err.Error(), nil)
	}

	return BuildSuccessResponse(req.ID, result)
}

// HasMethod returns true if a handler is registered for the method.
func (r *MessageRouter) HasMethod(method string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.handlers[method]
	_, existsCtx := r.handlersCtx[method]
	return exists || existsCtx
}

// Methods returns a list of all registered methods.
func (r *MessageRouter) Methods() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	methods := make([]string, 0, len(r.handlers)+len(r.handlersCtx))
	for method := range r.handlers {
		methods = append(methods, method)
	}
	for method := range r.handlersCtx {
		methods = append(methods, method)
	}
	return methods
}
