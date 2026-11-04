# ==============================================================================
# scripts/fixtures/apply-fixtures.ps1
# Applies all idempotent database test fixtures for multi-machine parity (Windows)
# ==============================================================================

param (
    [string]$DbHost = $(if ($env:DB_HOST) { $env:DB_HOST } else { "localhost" }),
    [string]$DbPort = $(if ($env:DB_PORT) { $env:DB_PORT } else { "5432" }),
    [string]$DbUser = $(if ($env:DB_USER) { $env:DB_USER } else { "postgres" }),
    [string]$DbName = $(if ($env:DB_NAME) { $env:DB_NAME } else { "dink_machinery" }),
    [string]$DbPassword = $(if ($env:DB_PASSWORD) { $env:DB_PASSWORD } else { "postgres" })
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host "==> Synchronizing Idempotent Database Test Fixtures [Target: ${DbName}@${DbHost}:${DbPort}]..." -ForegroundColor Cyan

$env:PGPASSWORD = $DbPassword

$fixtures = Get-ChildItem -Path $ScriptDir -Filter "*.sql" | Sort-Object Name

foreach ($file in $fixtures) {
    Write-Host "    --> Applying $($file.Name)..." -ForegroundColor Yellow
    & psql -h $DbHost -p $DbPort -U $DbUser -d $DbName -f $file.FullName | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to apply fixture: $($file.Name)"
        exit 1
    }
}

Write-Host "==> Database Test Fixtures Successfully Synchronized (Multi-Machine Parity Guaranteed)." -ForegroundColor Green
