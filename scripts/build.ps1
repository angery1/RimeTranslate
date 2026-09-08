$ErrorActionPreference = "Stop"

$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$dist = Join-Path $root "dist"
New-Item -ItemType Directory -Force -Path $dist | Out-Null

Push-Location $root
try {
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    gofmt -w .\src\main.go
    go vet .\src
    go build -trimpath -ldflags="-H=windowsgui -s -w" -o .\dist\RimeTranslate.exe .\src

    $hash = (Get-FileHash .\dist\RimeTranslate.exe -Algorithm SHA256).Hash
    Write-Host "Built: dist\RimeTranslate.exe" -ForegroundColor Green
    Write-Host "SHA256: $hash"
}
finally {
    Pop-Location
}
