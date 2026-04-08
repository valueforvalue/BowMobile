# Bow Release Build Script

Write-Host "--- Starting Release Build ---" -ForegroundColor Cyan

# 1. Generate Templates
Write-Host "Step 1: Generating Templ components..."
templ generate ./bow-gui/

# 2. Sync Data
# (No longer syncing DB into source for embedding)
if (-Not (Test-Path bow-gui/frontend/assets)) { mkdir bow-gui/frontend/assets }
cp assets/* bow-gui/frontend/assets/

# 3. Build Binary
Write-Host "Step 3: Building Bow_v1.1.1.exe..."
cd bow-gui
& "C:\Users\value\go\bin\wails.exe" build -o Bow_v1.1.1.exe -ldflags "-H windowsgui"
mv build/bin/Bow_v1.1.1.exe ../
cd ..

# 4. Package into Zip
Write-Host "Step 4: Packaging into Zip (EXE + Database)..."
if (Test-Path Bow_v1.1.1_Release.zip) { Remove-Item Bow_v1.1.1_Release.zip }
Compress-Archive -Path Bow_v1.1.1.exe, parts.db -DestinationPath Bow_v1.1.1_Release.zip

Write-Host "--- Build Complete: Bow_v1.1.1.exe and Bow_v1.1.1_Release.zip created ---" -ForegroundColor Green
Write-Host "Note: This version requires 'parts.db' to be in the same folder as the EXE." -ForegroundColor Yellow
