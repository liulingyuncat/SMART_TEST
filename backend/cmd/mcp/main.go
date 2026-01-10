// Package main provides the entry point for the MCP server.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"webtest/internal/mcp"
	"webtest/internal/mcp/tools/handlers"
)

const (
	// Version is the current version of the MCP server.
	Version = "1.0.0"
	// ServerName is the name of the MCP server.
	ServerName = "webtest-mcp-server"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("config", "./config/mcp-server.yaml", "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version and exit if requested
	if *showVersion {
		fmt.Printf("%s version %s\n", ServerName, Version)
		os.Exit(0)
	}

	// Create MCP server
	server, err := mcp.NewServer(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	// Register all tools
	handlers.RegisterAllTools(server.Registry(), server.BackendClient())
	logInfo("info", "Registered %d tools", server.Registry().Count())

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logInfo("info", "Received signal %v, initiating shutdown...", sig)
		cancel()
		server.Shutdown()
	}()

	// Run the server
	if err := server.Run(ctx); err != nil && err != context.Canceled {
		logInfo("error", "Server stopped with error: %v", err)
		os.Exit(1)
	}

	logInfo("info", "Server shutdown complete")
}

func logInfo(level string, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] %s\n", level, msg)
}
