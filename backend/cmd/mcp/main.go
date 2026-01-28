// Package main provides the entry point for the MCP server.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"webtest/internal/mcp"
	"webtest/internal/mcp/tools/handlers"

	"github.com/joho/godotenv"
)

const (
	// Version is the current version of the MCP server.
	Version = "1.0.0"
	// ServerName is the name of the MCP server.
	ServerName = "webtest-mcp-server"
)

func main() {
	// Load .env file from multiple possible locations
	loadEnvFile()

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

// loadEnvFile loads environment variables from .env file
// It searches for .env in multiple locations:
// 1. Current working directory
// 2. Parent directory (webtest root)
// 3. Two levels up (for backend/build execution)
func loadEnvFile() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logInfo("debug", "Could not get current directory: %v", err)
		return
	}

	// Try multiple possible .env locations
	possiblePaths := []string{
		filepath.Join(cwd, ".env"),             // Current directory
		filepath.Join(cwd, "..", ".env"),       // Parent directory
		filepath.Join(cwd, "..", "..", ".env"), // Two levels up
		filepath.Join(cwd, "backend", ".env"),  // From root to backend
		"D:\\VSCode\\webtest\\.env",            // Absolute path fallback
	}

	for _, envPath := range possiblePaths {
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err == nil {
				logInfo("info", "Loaded .env file from: %s", envPath)

				// Log PROMPTS_DIR if it was set
				if promptsDir := os.Getenv("PROMPTS_DIR"); promptsDir != "" {
					logInfo("info", "PROMPTS_DIR from .env: %s", promptsDir)

					// Convert relative path to absolute if needed
					if !filepath.IsAbs(promptsDir) {
						// Determine the backend directory based on .env location
						envDir := filepath.Dir(envPath)
						backendDir := filepath.Join(envDir, "backend")

						// If .env is in backend or backend/build, adjust accordingly
						if filepath.Base(envDir) == "backend" {
							backendDir = envDir
						} else if filepath.Base(envDir) == "build" {
							backendDir = filepath.Dir(envDir) // parent of build is backend
						}

						absPromptsDir := filepath.Join(backendDir, promptsDir)
						os.Setenv("PROMPTS_DIR", absPromptsDir)
						logInfo("info", "Converted PROMPTS_DIR to absolute: %s (backend dir: %s)", absPromptsDir, backendDir)
					}
				}
				return
			}
		}
	}

	logInfo("debug", "No .env file found in any of the searched locations")
}

func logInfo(level string, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] %s\n", level, msg)
}
