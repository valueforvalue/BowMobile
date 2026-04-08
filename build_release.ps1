# Bow Release Build Script

Write-Host "--- Starting Release Build ---" -ForegroundColor Cyan

# 1. Generate Templates
Write-Host "Step 1: Generating Templ components..."
templ generate ./cmd/server/

# 2. Sync Data for Embedding
Write-Host "Step 2: Syncing data for embedding..."
cp parts.db cmd/server/
if (-Not (Test-Path cmd/server/assets)) { mkdir cmd/server/assets }
cp assets/* cmd/server/assets/

# 3. Build Binary
# -ldflags="-H windowsgui" hides the console window on launch
Write-Host "Step 3: Building Bow.exe..."
go build -o Bow.exe -ldflags="-s -w -H windowsgui" ./cmd/server/

Write-Host "--- Build Complete: Bow.exe created ---" -ForegroundColor Green
Write-Host "Note: When you run Bow.exe, it will extract 'parts.db' to the current folder if it doesn't exist." -ForegroundColor Yellow
