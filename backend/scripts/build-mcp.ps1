# MCP Server Build and Test Script for Windows PowerShell

# Build the MCP server
function Build-MCPServer {
    param(
        [string]$OutputPath = "build/mcp-server.exe"
    )
    
    Write-Host "Building MCP Server..." -ForegroundColor Cyan
    
    # Ensure build directory exists
    $buildDir = Split-Path -Parent $OutputPath
    if (-not (Test-Path $buildDir)) {
        New-Item -ItemType Directory -Path $buildDir -Force | Out-Null
    }
    
    # Build
    $env:CGO_ENABLED = "0"
    go build -ldflags="-s -w" -o $OutputPath ./cmd/mcp/...
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build successful: $OutputPath" -ForegroundColor Green
        return $true
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
        return $false
    }
}

# Run MCP tests
function Test-MCP {
    param(
        [switch]$Verbose,
        [switch]$Coverage
    )
    
    Write-Host "Running MCP tests..." -ForegroundColor Cyan
    
    $args = @("test")
    
    if ($Verbose) {
        $args += "-v"
    }
    
    if ($Coverage) {
        $args += "-cover"
        $args += "-coverprofile=mcp_coverage.out"
    }
    
    $args += "./internal/mcp/..."
    
    & go @args
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "All tests passed!" -ForegroundColor Green
        
        if ($Coverage) {
            Write-Host "`nCoverage Report:" -ForegroundColor Cyan
            go tool cover -func=mcp_coverage.out
        }
        return $true
    } else {
        Write-Host "Some tests failed!" -ForegroundColor Red
        return $false
    }
}

# Clean build artifacts
function Clean-MCP {
    Write-Host "Cleaning build artifacts..." -ForegroundColor Cyan
    
    if (Test-Path "build/mcp-server.exe") {
        Remove-Item "build/mcp-server.exe" -Force
    }
    if (Test-Path "mcp_coverage.out") {
        Remove-Item "mcp_coverage.out" -Force
    }
    
    Write-Host "Clean complete." -ForegroundColor Green
}

# Main script logic
$command = $args[0]

switch ($command) {
    "build" {
        Build-MCPServer
    }
    "test" {
        Test-MCP -Verbose
    }
    "test-coverage" {
        Test-MCP -Verbose -Coverage
    }
    "clean" {
        Clean-MCP
    }
    "all" {
        if (Test-MCP) {
            Build-MCPServer
        }
    }
    default {
        Write-Host "MCP Server Build Script" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Usage: .\scripts\build-mcp.ps1 <command>" -ForegroundColor White
        Write-Host ""
        Write-Host "Commands:" -ForegroundColor Yellow
        Write-Host "  build          Build the MCP server"
        Write-Host "  test           Run tests"
        Write-Host "  test-coverage  Run tests with coverage"
        Write-Host "  clean          Clean build artifacts"
        Write-Host "  all            Run tests then build"
    }
}
