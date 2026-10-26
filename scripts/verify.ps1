param([string]$Go = "go")
$ErrorActionPreference = "Stop"
$repo = Split-Path $PSScriptRoot -Parent
Push-Location (Join-Path $repo "combined_refactor")
$oldGOOS = $env:GOOS
$oldGOARCH = $env:GOARCH
$oldCGO = $env:CGO_ENABLED
try {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"
    & $Go test ./...
    if ($LASTEXITCODE -ne 0) { throw "Go tests failed" }
    & $Go vet ./...
    if ($LASTEXITCODE -ne 0) { throw "Go vet failed" }
    New-Item -ItemType Directory -Force (Join-Path $repo "release_assets") | Out-Null
    & $Go build -trimpath -o ../release_assets/cfdata-test.exe .
    if ($LASTEXITCODE -ne 0) { throw "Windows build failed" }
    $env:GOOS = "linux"
    & $Go build -trimpath -ldflags="-s -w" -o ../release_assets/cfdata-test .
    if ($LASTEXITCODE -ne 0) { throw "Linux build failed" }
    $env:GOOS = "android"
    $env:GOARCH = "arm64"
    & $Go build -trimpath -o ../release_assets/cfdata-test-android-arm64 .
    if ($LASTEXITCODE -ne 0) { throw "Android backend build failed" }
} finally {
    $env:GOOS = $oldGOOS
    $env:GOARCH = $oldGOARCH
    $env:CGO_ENABLED = $oldCGO
    Pop-Location
}
