# SMART_TEST - Server Startup Script
# Auto-load .env and start backend server

$ErrorActionPreference = "Stop"

# Project paths
$ProjectRoot = Split-Path -Parent $PSScriptRoot
$EnvFile = Join-Path $ProjectRoot ".env"
$BackendDir = Join-Path $ProjectRoot "backend"

Write-Host "[1/3] Loading environment variables..." -ForegroundColor Cyan

# Read .env file and set environment variables
if (Test-Path $EnvFile) {
    Get-Content $EnvFile | ForEach-Object {
        $line = $_.Trim()
        if ($line -and -not $line.StartsWith("#")) {
            $parts = $line -split "=", 2
            if ($parts.Length -eq 2) {
                $key = $parts[0].Trim()
                $value = $parts[1].Trim()
                [Environment]::SetEnvironmentVariable($key, $value, "Process")
                Write-Host "  Set: $key = $value" -ForegroundColor Green
            }
        }
    }
} else {
    Write-Host "  Warning: .env file not found, using defaults" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "[2/3] Building backend..." -ForegroundColor Cyan
Set-Location $BackendDir
go build -o server.exe ./cmd/server
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "  Build complete" -ForegroundColor Green

Write-Host ""
Write-Host "[3/3] Starting server..." -ForegroundColor Cyan
Write-Host "  URL: https://localhost:8443" -ForegroundColor White
Write-Host ""

.\server.exe
