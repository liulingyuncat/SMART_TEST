package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// DebugMode controls whether debug logs are printed (disabled by default for performance)
var DebugMode = os.Getenv("MCP_DEBUG") == "1"

// HTTPTransport implements Transport using HTTP for MCP communication.
// It implements the Streamable HTTP transport as defined in MCP spec.
type HTTPTransport struct {
	server         *http.Server
	addr           string
	path           string
	msgCh          chan *Message
	respCh         chan []byte
	closeCh        chan struct{}
	logger         *log.Logger
	mu             sync.Mutex
	closed         bool
	currentReq     chan []byte    // Channel for current request's response
	requestHandler RequestHandler // Direct request handler for sync mode
}

// debugLog only logs when DebugMode is enabled
func (t *HTTPTransport) debugLog(format string, args ...interface{}) {
	if DebugMode {
		t.logger.Printf(format, args...)
	}
}

// NewHTTPTransport creates a new HTTP transport on the specified address.
func NewHTTPTransport(addr string, path string) *HTTPTransport {
	if path == "" {
		path = "/mcp"
	}
	return &HTTPTransport{
		addr:    addr,
		path:    path,
		msgCh:   make(chan *Message, 100),
		respCh:  make(chan []byte, 100),
		closeCh: make(chan struct{}),
		logger:  log.New(os.Stderr, "[HTTP] ", log.LstdFlags),
	}
}

// SetRequestHandler sets the handler for synchronous request processing.
func (t *HTTPTransport) SetRequestHandler(handler RequestHandler) {
	t.requestHandler = handler
}

// Start initializes and starts the HTTP server.
func (t *HTTPTransport) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc(t.path, t.handleMCP)

	t.server = &http.Server{
		Addr:         t.addr,
		Handler:      mux,
		ReadTimeout:  0,                // No timeout - allow long-running requests
		WriteTimeout: 0,                // No timeout - allow long-running responses
		IdleTimeout:  10 * time.Minute, // Extended idle timeout (VS Code MCP client may wait long)
	}

	// Create listener first to catch port binding errors early
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		t.logger.Printf("ERROR: Failed to bind to %s: %v", t.addr, err)
		return fmt.Errorf("failed to bind to %s: %w", t.addr, err)
	}

	// Wrap with TCP keep-alive listener
	tcpListener := listener.(*net.TCPListener)
	keepAliveListener := tcpKeepAliveListener{tcpListener}

	t.logger.Printf("INFO: HTTP server bound to %s", listener.Addr().String())
	t.logger.Printf("INFO: MCP endpoint available at http://localhost%s%s", t.addr, t.path)
	t.debugLog("DEBUG: Server configuration - ReadTimeout: unlimited, WriteTimeout: unlimited, IdleTimeout: 10 minutes, TCP Keep-Alive: 30s")

	// Enable HTTP keep-alives
	t.server.SetKeepAlivesEnabled(true)

	go func() {
		t.logger.Printf("INFO: Starting HTTP server, listening for connections...")
		if err := t.server.Serve(keepAliveListener); err != nil && err != http.ErrServerClosed {
			t.logger.Printf("ERROR: HTTP server error: %v", err)
		} else {
			t.logger.Printf("INFO: HTTP server stopped")
		}
	}()

	// Wait for context cancellation
	go func() {
		<-ctx.Done()
		t.logger.Printf("INFO: Context cancelled, initiating shutdown...")
		t.Close()
	}()

	return nil
}

// handleMCP handles MCP requests over HTTP.
func (t *HTTPTransport) handleMCP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	requestID := fmt.Sprintf("%d", startTime.UnixNano())

	// Log incoming request details (only in debug mode)
	t.debugLog("DEBUG: [%s] Incoming request from %s", requestID, r.RemoteAddr)
	t.debugLog("DEBUG: [%s] Method: %s, URL: %s, Proto: %s", requestID, r.Method, r.URL.String(), r.Proto)
	t.debugLog("DEBUG: [%s] Headers: %v", requestID, t.sanitizeHeaders(r.Header))

	// Set CORS headers for VS Code MCP client
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Token, Accept")
	w.Header().Set("Access-Control-Max-Age", "86400")
	// Enable HTTP Keep-Alive for persistent connections
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Keep-Alive", "timeout=600, max=1000")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only accept POST requests
	if r.Method != http.MethodPost {
		t.logger.Printf("WARN: [%s] Method not allowed: %s", requestID, r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.logger.Printf("ERROR: [%s] Error reading request body: %v", requestID, err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		t.logger.Printf("WARN: [%s] Empty request body", requestID)
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// Log request body only in debug mode (truncated for large bodies)
	if DebugMode {
		bodyStr := string(body)
		if len(bodyStr) > 500 {
			t.debugLog("DEBUG: [%s] Request body (truncated): %s...", requestID, bodyStr[:500])
		} else {
			t.debugLog("DEBUG: [%s] Request body: %s", requestID, bodyStr)
		}
	}

	// Try to parse method for logging and detect notifications
	var req MCPRequest
	isNotification := false
	if err := json.Unmarshal(body, &req); err == nil {
		t.debugLog("DEBUG: [%s] MCP method: %s", requestID, req.Method)
		// A notification has no "id" field (ID is nil or empty)
		isNotification = len(req.ID) == 0 || string(req.ID) == "null"
	}

	// Extract authentication token from headers
	metadata := make(map[string]string)
	if token := r.Header.Get("X-API-Token"); token != "" {
		metadata["api_token"] = token
	} else if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		// Extract Bearer token
		if strings.HasPrefix(authHeader, "Bearer ") {
			metadata["api_token"] = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// For notifications, return immediately with 202 Accepted (no response expected)
	if isNotification {
		// Still queue the notification for async processing
		msg := &Message{Data: body, Metadata: metadata}
		select {
		case t.msgCh <- msg:
		default:
			// Drop notification if queue is full
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}

	// SYNC MODE: If we have a request handler, process directly without queue
	if t.requestHandler != nil {
		// Create request context with timeout
		reqCtx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Process request synchronously
		resp := t.requestHandler(reqCtx, body, metadata)
		if resp != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
			t.debugLog("DEBUG: [%s] Sync request completed in %v", requestID, time.Since(startTime))
		}
		return
	}

	// ASYNC MODE (fallback): Use queue-based processing
	respCh := make(chan []byte, 1)
	t.mu.Lock()
	t.currentReq = respCh
	t.mu.Unlock()

	msg := &Message{Data: body, Metadata: metadata}
	select {
	case t.msgCh <- msg:
	case <-t.closeCh:
		t.logger.Printf("ERROR: [%s] Server shutting down, rejecting request", requestID)
		http.Error(w, "Server shutting down", http.StatusServiceUnavailable)
		return
	}

	select {
	case resp := <-respCh:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		t.debugLog("DEBUG: [%s] Async request completed in %v", requestID, time.Since(startTime))
	case <-t.closeCh:
		t.logger.Printf("ERROR: [%s] Server shutting down during request", requestID)
		http.Error(w, "Server shutting down", http.StatusServiceUnavailable)
	}
}

// sanitizeHeaders creates a safe string representation of headers for logging
func (t *HTTPTransport) sanitizeHeaders(headers http.Header) string {
	safe := make(map[string]string)
	for k, v := range headers {
		if strings.ToLower(k) == "authorization" || strings.ToLower(k) == "x-api-token" {
			safe[k] = "[REDACTED]"
		} else {
			safe[k] = strings.Join(v, ", ")
		}
	}
	b, _ := json.Marshal(safe)
	return string(b)
}

// Send sends a response message back to the HTTP client.
func (t *HTTPTransport) Send(message []byte) error {
	t.mu.Lock()
	respCh := t.currentReq
	t.mu.Unlock()

	if respCh != nil {
		select {
		case respCh <- message:
			return nil
		case <-t.closeCh:
			t.logger.Printf("ERROR: Transport closed while sending response")
			return fmt.Errorf("transport is closed")
		default:
			t.logger.Printf("ERROR: Response channel is full")
			return fmt.Errorf("response channel full")
		}
	}
	t.logger.Printf("ERROR: No active request to send response to")
	return fmt.Errorf("no active request")
}

// Receive waits for and returns the next message from HTTP requests.
func (t *HTTPTransport) Receive() ([]byte, error) {
	msg, err := t.ReceiveWithMetadata()
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

// ReceiveWithMetadata waits for and returns the next message with metadata.
func (t *HTTPTransport) ReceiveWithMetadata() (*Message, error) {
	select {
	case msg := <-t.msgCh:
		return msg, nil
	case <-t.closeCh:
		t.logger.Printf("INFO: Transport closed, returning EOF")
		return nil, io.EOF
	}
}

// Close shuts down the HTTP transport.
func (t *HTTPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.logger.Printf("INFO: Closing HTTP transport...")
	t.closed = true
	close(t.closeCh)

	if t.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t.logger.Printf("INFO: Shutting down HTTP server with 5s timeout")
		err := t.server.Shutdown(ctx)
		if err != nil {
			t.logger.Printf("ERROR: Error during server shutdown: %v", err)
		} else {
			t.logger.Printf("INFO: HTTP server shutdown complete")
		}
		return err
	}

	return nil
}

// MCPRequest represents a JSON-RPC request for logging purposes.
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
}

// Log writes a message to stderr for debugging purposes.
func (t *HTTPTransport) Log(format string, args ...interface{}) {
	t.logger.Printf(format, args...)
}

// tcpKeepAliveListener wraps a net.TCPListener to enable TCP keep-alive on accepted connections
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts a connection and enables TCP keep-alive
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	// Enable TCP keep-alive with 30 second probe interval
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(30 * time.Second)
	return tc, nil
}
