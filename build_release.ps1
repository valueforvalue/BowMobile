# Bow Release Build Script

Write-Host "--- Starting Release Build ---" -ForegroundColor Cyan

# 0. Deep Clean
Write-Host "Step 0: Cleaning old build artifacts..."
Remove-Item -Path *.exe -ErrorAction SilentlyContinue
Remove-Item -Path *.zip -ErrorAction SilentlyContinue
if (Test-Path bow-gui/build/bin) { Remove-Item -Path bow-gui/build/bin/* -Recurse -Force -ErrorAction SilentlyContinue }

# 1. Generate Templates
Write-Host "Step 1: Generating Templ components..."
templ generate ./bow-gui/
templ generate ./cmd/server/

# 2. Sync Data
# (No longer syncing DB into source for embedding)
if (-Not (Test-Path bow-gui/frontend/assets)) { mkdir bow-gui/frontend/assets }
cp assets/* bow-gui/frontend/assets/

# 2b. Verify Search Logic
Write-Host "Step 2b: Verifying search logic..."
cd cmd/tools
go run verify_search.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Verification FAILED. Aborting build." -ForegroundColor Red
    cd ../..
    exit 1
}
cd ../..

# 3. Build Binary
Write-Host "Step 3: Building Bow_v1.1.3.exe..."
cd bow-gui
wails build -o Bow_v1.1.3.exe -ldflags "-H windowsgui"
mv build/bin/Bow_v1.1.3.exe ../
cd ..

# 4. Package into Zip
Write-Host "Step 4: Packaging into Zip (EXE + Database)..."
if (Test-Path Bow_v1.1.3_Release.zip) { Remove-Item Bow_v1.1.3_Release.zip }
Compress-Archive -Path Bow_v1.1.3.exe, parts.db -DestinationPath Bow_v1.1.3_Release.zip

Write-Host "--- Build Complete: Bow_v1.1.3.exe and Bow_v1.1.3_Release.zip created ---" -ForegroundColor Green
Write-Host "Note: This version requires 'parts.db' to be in the same folder as the EXE." -ForegroundColor Yellow
