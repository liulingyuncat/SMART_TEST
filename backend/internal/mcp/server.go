// Package mcp provides the MCP server implementation.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"webtest/internal/mcp/client"
	"webtest/internal/mcp/config"
	"webtest/internal/mcp/prompts"
	"webtest/internal/mcp/protocol"
	"webtest/internal/mcp/tools"
	"webtest/internal/mcp/transport"
)

const (
	// ServerName is the name of the MCP server.
	ServerName = "webtest-mcp-server"
	// ServerVersion is the version of the MCP server.
	ServerVersion = "1.0.0"
	// ProtocolVersion is the MCP protocol version supported.
	ProtocolVersion = "2025-06-18"
)

// Server represents the MCP server.
type Server struct {
	config          *config.Config
	transport       transport.Transport
	router          *protocol.MessageRouter
	registry        *tools.ToolRegistry
	promptsRegistry *prompts.PromptsRegistry
	backendClient   *client.BackendClient
	validator       *tools.SchemaValidator
	logger          *log.Logger
	initialized     bool
}

// NewServer creates a new MCP server with the given configuration path.
func NewServer(configPath string) (*Server, error) {
	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create logger (output to stderr to not interfere with stdio transport)
	logger := log.New(os.Stderr, "[MCP] ", log.LstdFlags)

	// Create auth manager
	authManager, err := client.NewAuthManager(cfg.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth manager: %w", err)
	}

	// Log authentication mode
	if cfg.Auth.DynamicToken {
		logger.Println("Dynamic token mode enabled - token will be read from request headers")
	} else {
		// Validate token if configured and not using dynamic token
		if cfg.Auth.ValidateOnStart {
			logger.Println("Validating authentication token...")
			if err := authManager.ValidateToken(context.Background(), cfg.Backend.BaseURL); err != nil {
				return nil, fmt.Errorf("token validation failed: %w", err)
			}
			logger.Println("Token validated successfully")
		}
	}

	// Create backend client
	backendClient := client.NewBackendClient(cfg.Backend, authManager)

	// Create components
	registry := tools.NewToolRegistry()
	validator := tools.NewSchemaValidator()
	router := protocol.NewMessageRouter()

	// Create prompts registry and load system prompts
	promptsRegistry := prompts.NewPromptsRegistry()
	promptLoader := &prompts.PromptLoader{}

	// Log current working directory
	cwd, _ := os.Getwd()
	logger.Printf("DEBUG: Current working directory: %s", cwd)

	// Get the directory of the executable
	exePath, err := os.Executable()
	logger.Printf("DEBUG: Executable path: %s (error: %v)", exePath, err)

	var promptsDir string

	// Try multiple possible paths
	possiblePaths := []string{}

	if err == nil {
		exeDir := filepath.Dir(exePath)
		// The executable is in backend/ directory, so prompts are in backend/internal/mcp/prompts/
		path1 := filepath.Join(exeDir, "internal", "mcp", "prompts")
		possiblePaths = append(possiblePaths, path1)
		logger.Printf("DEBUG: Possible path 1 (relative to exe): %s", path1)
	}

	// Try relative to cwd
	path2 := filepath.Join(cwd, "internal", "mcp", "prompts")
	possiblePaths = append(possiblePaths, path2)
	logger.Printf("DEBUG: Possible path 2 (relative to cwd): %s", path2)

	// Try from project root
	path3 := filepath.Join("D:\\VSCode\\webtest\\backend", "internal", "mcp", "prompts")
	possiblePaths = append(possiblePaths, path3)
	logger.Printf("DEBUG: Possible path 3 (absolute): %s", path3)

	// Find the first valid path
	for _, p := range possiblePaths {
		if stat, err := os.Stat(p); err == nil && stat.IsDir() {
			promptsDir = p
			logger.Printf("DEBUG: Found valid prompts directory: %s", promptsDir)
			break
		}
	}

	if promptsDir == "" {
		logger.Printf("ERROR: Could not find prompts directory in any of the possible paths")
		for _, p := range possiblePaths {
			if stat, err := os.Stat(p); err != nil {
				logger.Printf("  - %s (not found: %v)", p, err)
			} else {
				logger.Printf("  - %s (exists: isDir=%v)", p, stat.IsDir())
			}
		}
	} else {
		logger.Printf("DEBUG: Using prompts directory: %s", promptsDir)
		// List files in directory
		if entries, err := os.ReadDir(promptsDir); err == nil {
			logger.Printf("DEBUG: Directory contains %d entries:", len(entries))
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".prompt.md") {
					logger.Printf("  - %s (IsDir: %v)", entry.Name(), entry.IsDir())
				}
			}
		}
	}

	if promptsDir != "" {
		if err := promptLoader.LoadAll(promptsDir, promptsRegistry); err != nil {
			logger.Printf("WARNING: Failed to load system prompts from %s: %v", promptsDir, err)
		} else {
			count := promptsRegistry.Count()
			logger.Printf("INFO: Loaded %d system prompts from %s", count, promptsDir)
			if count == 0 {
				logger.Printf("WARNING: No prompts were loaded, the directory might be empty")
			}
		}
	}

	// Create transport based on mode
	var trans transport.Transport
	switch cfg.Server.Mode {
	case "stdio":
		trans = transport.NewStdioTransport()
	case "http":
		trans = transport.NewHTTPTransport(cfg.Server.HTTPAddr, cfg.Server.HTTPPath)
	default:
		return nil, fmt.Errorf("unsupported transport mode: %s", cfg.Server.Mode)
	}

	server := &Server{
		config:          cfg,
		transport:       trans,
		router:          router,
		registry:        registry,
		promptsRegistry: promptsRegistry,
		backendClient:   backendClient,
		validator:       validator,
		logger:          logger,
		initialized:     false,
	}

	// Register protocol handlers
	server.registerHandlers()

	return server, nil
}

// registerHandlers registers the MCP protocol handlers.
func (s *Server) registerHandlers() {
	s.router.Register("initialize", s.handleInitialize)
	s.router.Register("initialized", s.handleInitialized)
	s.router.Register("tools/list", s.handleToolsList)
	s.router.RegisterWithContext("prompts/list", s.handlePromptsList)
	s.router.RegisterWithContext("prompts/get", s.handlePromptsGet)
	s.router.RegisterWithContext("tools/call", s.handleToolsCallWithContext)
}

// Run starts the MCP server main loop.
func (s *Server) Run(ctx context.Context) error {
	s.logger.Printf("INFO: Starting %s v%s", ServerName, ServerVersion)
	s.logger.Printf("INFO: Transport mode: %s", s.config.Server.Mode)
	s.logger.Printf("INFO: Registered prompts: %d", s.promptsRegistry.Count())
	s.logger.Printf("INFO: Backend URL: %s", s.config.Backend.BaseURL)
	s.logger.Printf("INFO: Registered tools: %d", s.registry.Count())
	s.logger.Printf("INFO: Dynamic token mode: %v", s.config.Auth.DynamicToken)

	// Set up synchronous request handler for HTTP transport
	if syncTransport, ok := s.transport.(transport.SyncTransport); ok {
		s.logger.Printf("INFO: Enabling synchronous request handling for better performance")
		syncTransport.SetRequestHandler(s.handleSyncRequest)
	}

	// Start transport
	s.logger.Printf("INFO: Starting transport...")
	if err := s.transport.Start(ctx); err != nil {
		s.logger.Printf("ERROR: Failed to start transport: %v", err)
		return fmt.Errorf("failed to start transport: %w", err)
	}
	defer func() {
		s.logger.Printf("INFO: Closing transport...")
		s.transport.Close()
	}()

	s.logger.Printf("INFO: Server ready, waiting for requests...")

	// Main loop with fault tolerance
	errorCount := 0
	maxConsecutiveErrors := 10

	for {
		select {
		case <-ctx.Done():
			s.logger.Printf("INFO: Context cancelled, shutting down...")
			return ctx.Err()
		default:
			msg, err := s.transport.ReceiveWithMetadata()
			if err != nil {
				if err == io.EOF {
					s.logger.Printf("INFO: Connection closed (EOF)")
					return nil
				}
				errorCount++
				s.logger.Printf("WARNING: Error receiving message (attempt %d/%d): %v", errorCount, maxConsecutiveErrors, err)

				// Only exit if we have too many consecutive errors
				if errorCount > maxConsecutiveErrors {
					s.logger.Printf("ERROR: Too many consecutive errors (%d), shutting down", errorCount)
					return fmt.Errorf("too many consecutive receive errors: %w", err)
				}
				// Continue processing instead of stopping
				continue
			}

			// Reset error count on successful receive
			errorCount = 0

			s.logger.Printf("DEBUG: Processing message, length: %d bytes", len(msg.Data))

			// Create context with token from metadata if available
			reqCtx := ctx
			if token, ok := msg.Metadata["api_token"]; ok && token != "" {
				s.logger.Printf("DEBUG: Using token from request headers")
				reqCtx = context.WithValue(ctx, client.TokenContextKey, token)
			}

			// Handle message with panic recovery
			func() {
				defer func() {
					if r := recover(); r != nil {
						s.logger.Printf("ERROR: Panic recovered in message handler: %v", r)
					}
				}()

				response := s.handleMessageWithContext(reqCtx, msg.Data)
				if response != nil {
					respBytes, err := response.Serialize()
					if err != nil {
						s.logger.Printf("ERROR: Error serializing response: %v", err)
						return
					}
					s.logger.Printf("DEBUG: Sending response, length: %d bytes", len(respBytes))
					if err := s.transport.Send(respBytes); err != nil {
						s.logger.Printf("WARNING: Error sending response: %v (connection may be closed)", err)
						// Don't exit on send error - client may have disconnected
					} else {
						s.logger.Printf("DEBUG: Response sent successfully")
					}
				}
			}()
		}
	}
}

// handleSyncRequest processes a request synchronously (used by HTTP transport).
// This is called directly from the HTTP handler goroutine for low latency.
func (s *Server) handleSyncRequest(ctx context.Context, data []byte, metadata map[string]string) []byte {
	// Add token to context if available
	if token, ok := metadata["api_token"]; ok && token != "" {
		ctx = context.WithValue(ctx, client.TokenContextKey, token)
	}

	// Process the request
	response := s.handleMessageWithContext(ctx, data)
	if response == nil {
		return nil
	}

	// Serialize response
	respBytes, err := response.Serialize()
	if err != nil {
		s.logger.Printf("ERROR: Error serializing response: %v", err)
		return nil
	}
	return respBytes
}

// handleMessage processes a single incoming message.
func (s *Server) handleMessage(data []byte) *protocol.Response {
	return s.handleMessageWithContext(context.Background(), data)
}

// handleMessageWithContext processes a single incoming message with context.
func (s *Server) handleMessageWithContext(ctx context.Context, data []byte) *protocol.Response {
	// Recover from panics in message handling
	defer func() {
		if r := recover(); r != nil {
			s.logger.Printf("ERROR: Panic in handleMessageWithContext: %v", r)
		}
	}()

	// Parse request
	req, err := protocol.ParseRequest(data)
	if err != nil {
		s.logger.Printf("ERROR: Failed to parse request: %v", err)
		if protoErr, ok := err.(*protocol.Error); ok {
			return protocol.BuildErrorResponse(nil, protoErr.Code, protoErr.Message, protoErr.Data)
		}
		return protocol.BuildErrorResponse(nil, protocol.ParseError, "Parse error", err.Error())
	}

	s.logger.Printf("INFO: Received method: %s (ID: %v, Notification: %v)", req.Method, req.ID, req.IsNotification())

	// Store context in request for handlers to access
	req.Context = ctx

	// Route to handler with error handling
	response := func() *protocol.Response {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Printf("ERROR: Panic in route handler for method %s: %v", req.Method, r)
			}
		}()
		return s.router.Route(req)
	}()

	// Don't send response for notifications
	if req.IsNotification() {
		s.logger.Printf("DEBUG: Notification processed, no response needed")
		return nil
	}

	s.logger.Printf("DEBUG: Handler returned response for method: %s", req.Method)
	return response
}

// handleInitialize handles the initialize request.
func (s *Server) handleInitialize(params json.RawMessage) (interface{}, error) {
	// Parse initialize params (optional, for logging)
	var initParams struct {
		ProtocolVersion string `json:"protocolVersion"`
		ClientInfo      struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"clientInfo"`
	}
	if err := json.Unmarshal(params, &initParams); err == nil {
		s.logger.Printf("Client: %s v%s, Protocol: %s",
			initParams.ClientInfo.Name,
			initParams.ClientInfo.Version,
			initParams.ProtocolVersion)
	}

	// Return server info
	return map[string]interface{}{
		"protocolVersion": ProtocolVersion,
		"capabilities": map[string]interface{}{
			"tools":   map[string]interface{}{},
			"prompts": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    ServerName,
			"version": ServerVersion,
		},
	}, nil
}

// handleInitialized handles the initialized notification.
func (s *Server) handleInitialized(params json.RawMessage) (interface{}, error) {
	s.initialized = true
	s.logger.Println("Client initialized")
	return nil, nil
}

// handleToolsList handles the tools/list request.
func (s *Server) handleToolsList(params json.RawMessage) (interface{}, error) {
	toolDefs := s.registry.List()
	s.logger.Printf("Returning %d tools", len(toolDefs))
	return map[string]interface{}{
		"tools": toolDefs,
	}, nil
}

// handleToolsCallWithContext handles the tools/call request with context.
func (s *Server) handleToolsCallWithContext(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// Parse call params
	var callParams struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	if err := json.Unmarshal(params, &callParams); err != nil {
		s.logger.Printf("ERROR: Failed to parse tool call params: %v", err)
		return nil, &protocol.Error{
			Code:    protocol.InvalidParams,
			Message: "Invalid params",
			Data:    err.Error(),
		}
	}

	s.logger.Printf("INFO: Tool call request - name: %s", callParams.Name)
	s.logger.Printf("DEBUG: Tool arguments: %v", callParams.Arguments)

	// Find tool
	handler := s.registry.Get(callParams.Name)
	if handler == nil {
		s.logger.Printf("ERROR: Unknown tool: %s", callParams.Name)
		return nil, &protocol.Error{
			Code:    protocol.InvalidParams,
			Message: fmt.Sprintf("Unknown tool: %s", callParams.Name),
		}
	}

	s.logger.Printf("DEBUG: Found handler for tool: %s", callParams.Name)

	// Validate arguments
	if err := s.validator.Validate(handler.InputSchema(), callParams.Arguments); err != nil {
		s.logger.Printf("ERROR: Argument validation failed for tool %s: %v", callParams.Name, err)
		if valErr, ok := err.(*tools.ValidationError); ok {
			return nil, &protocol.Error{
				Code:    protocol.InvalidParams,
				Message: "Invalid params",
				Data: map[string]interface{}{
					"field":  valErr.Field,
					"reason": valErr.Reason,
				},
			}
		}
		return nil, &protocol.Error{
			Code:    protocol.InvalidParams,
			Message: err.Error(),
		}
	}

	s.logger.Printf("DEBUG: Arguments validated successfully for tool: %s", callParams.Name)
	s.logger.Printf("INFO: Executing tool: %s", callParams.Name)

	// Execute tool with panic recovery
	var result tools.ToolResult
	var execErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Printf("ERROR: Panic during tool execution of %s: %v", callParams.Name, r)
				execErr = fmt.Errorf("tool execution panic: %v", r)
			}
		}()
		result, execErr = handler.Execute(ctx, callParams.Arguments)
	}()

	if execErr != nil {
		s.logger.Printf("ERROR: Tool %s execution failed: %v", callParams.Name, execErr)
		if protoErr, ok := execErr.(*protocol.Error); ok {
			return nil, protoErr
		}
		return nil, &protocol.Error{
			Code:    protocol.InternalError,
			Message: execErr.Error(),
		}
	}

	s.logger.Printf("INFO: Tool %s executed successfully", callParams.Name)
	if resultBytes, err := json.Marshal(result); err == nil {
		resultStr := string(resultBytes)
		if len(resultStr) > 200 {
			s.logger.Printf("DEBUG: Tool result (truncated): %s...", resultStr[:200])
		} else {
			s.logger.Printf("DEBUG: Tool result: %s", resultStr)
		}
	}

	// 验证工具结果的有效性（result是tools.ToolResult结构体）
	if len(result.Content) == 0 {
		s.logger.Printf("WARNING: Tool %s returned empty content, returning default response", callParams.Name)
		// Return a default response instead of erroring out
		return tools.NewTextResult("Tool completed but returned no content"), nil
	}
	s.logger.Printf("DEBUG: Tool %s result validation passed, content items: %d", callParams.Name, len(result.Content))

	// 尝试序列化验证
	if _, err := json.Marshal(result); err != nil {
		s.logger.Printf("ERROR: Tool %s result serialization failed: %v", callParams.Name, err)
		return nil, &protocol.Error{
			Code:    protocol.InternalError,
			Message: "Tool result serialization failed",
			Data:    err.Error(),
		}
	}

	return result, nil
}

// handleToolsCall handles the tools/call request (deprecated, use handleToolsCallWithContext).
func (s *Server) handleToolsCall(params json.RawMessage) (interface{}, error) {
	return s.handleToolsCallWithContext(context.Background(), params)
}

// Registry returns the tool registry for registering tools.
func (s *Server) Registry() *tools.ToolRegistry {
	return s.registry
}

// BackendClient returns the backend client for tool handlers.
func (s *Server) BackendClient() *client.BackendClient {
	return s.backendClient
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown() error {
	s.logger.Println("Shutting down...")
	return s.transport.Close()
}

// handlePromptsList handles the prompts/list request.
func (s *Server) handlePromptsList(ctx context.Context, params json.RawMessage) (interface{}, error) {
	s.logger.Printf("INFO: Handling prompts/list request")
	result, err := prompts.HandlePromptsList(ctx, s.promptsRegistry, s.backendClient, params)
	if err != nil {
		s.logger.Printf("ERROR: prompts/list failed: %v", err)
		return nil, err
	}
	s.logger.Printf("DEBUG: prompts/list returned successfully")
	return result, nil
}

// handlePromptsGet handles the prompts/get request.
func (s *Server) handlePromptsGet(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var getParams struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(params, &getParams); err == nil {
		s.logger.Printf("INFO: Handling prompts/get request for: %s", getParams.Name)
	}
	result, err := prompts.HandlePromptsGet(ctx, s.promptsRegistry, s.backendClient, params)
	if err != nil {
		s.logger.Printf("ERROR: prompts/get failed: %v", err)
		return nil, &protocol.Error{
			Code:    protocol.InvalidParams,
			Message: err.Error(),
		}
	}
	s.logger.Printf("DEBUG: prompts/get returned successfully")
	return result, nil
}
