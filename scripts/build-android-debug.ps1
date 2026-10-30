param([string]$Go = "go", [string]$Gradle = "gradle")
$ErrorActionPreference = "Stop"
$repo = Split-Path $PSScriptRoot -Parent
$oldGOOS = $env:GOOS
$oldGOARCH = $env:GOARCH
$oldCGO = $env:CGO_ENABLED
Push-Location $repo
try {
    & java -version
    if ($LASTEXITCODE -ne 0) { throw "Java 17 is required" }
    & $Gradle --version
    if ($LASTEXITCODE -ne 0) { throw "Gradle 8.10.2 is required" }
    New-Item -ItemType Directory -Force app/src/main/jniLibs/arm64-v8a | Out-Null
    $env:GOOS = "android"
    $env:GOARCH = "arm64"
    $env:CGO_ENABLED = "0"
    $versionLine = Get-Content app/version.properties | Where-Object { $_ -match '^versionName=' }
    $backendVersion = ($versionLine -split '=', 2)[1]
    if (!$backendVersion) { throw "Missing backend version" }
    Push-Location combined_refactor
    try {
        & $Go build -trimpath "-ldflags=-X main.appVersion=$backendVersion" -o ../app/src/main/jniLibs/arm64-v8a/libcfdata.so .
        if ($LASTEXITCODE -ne 0) { throw "Android backend build failed" }
    } finally { Pop-Location }
    & $Gradle assembleDebug lintDebug --no-daemon
    if ($LASTEXITCODE -ne 0) { throw "APK build/lint failed" }
    Get-FileHash app/build/outputs/apk/debug/app-debug.apk -Algorithm SHA256
} finally {
    $env:GOOS = $oldGOOS
    $env:GOARCH = $oldGOARCH
    $env:CGO_ENABLED = $oldCGO
    Pop-Location
}
