$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Release package: install.ps1 sits beside RimeTranslate.exe and rime\.
# Source repository: install.ps1 lives in scripts\, so project root is one level up.
if ((Test-Path (Join-Path $scriptDir "RimeTranslate.exe")) -and (Test-Path (Join-Path $scriptDir "rime"))) {
    $root = $scriptDir
} else {
    $root = Split-Path -Parent $scriptDir
}

$exeCandidates = @(
    (Join-Path $root "RimeTranslate.exe"),
    (Join-Path $root "dist\RimeTranslate.exe")
)
$exe = $exeCandidates | Where-Object { Test-Path $_ } | Select-Object -First 1
if (-not $exe) {
    throw "RimeTranslate.exe not found. Use a Release package or run scripts\build.ps1 first."
}

$rime = Join-Path $env:APPDATA "Rime"
$lua = Join-Path $rime "lua"
$backup = Join-Path $rime ("backup_RimeTranslate_" + (Get-Date -Format "yyyyMMdd_HHmmss"))

New-Item -ItemType Directory -Force -Path $lua | Out-Null
New-Item -ItemType Directory -Force -Path $backup | Out-Null

foreach ($name in @("async_ollama_filter.lua", "async_refresh.lua")) {
    $src = Join-Path $root ("rime\lua\" + $name)
    $dst = Join-Path $lua $name
    if (Test-Path $dst) {
        Copy-Item $dst (Join-Path $backup $name) -Force
    }
    Copy-Item $src $dst -Force
}

$runKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
New-Item -Path $runKey -Force | Out-Null
Set-ItemProperty -Path $runKey -Name "RimeTranslate" -Value ('"' + $exe + '"')

Get-Process -Name "RimeTranslate" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 300
Start-Process -FilePath $exe

Write-Host "" 
Write-Host "RimeTranslate installed." -ForegroundColor Green
Write-Host "Lua backup: $backup"
Write-Host "Startup: enabled for current user"
Write-Host ""
Write-Host "NEXT STEP:" -ForegroundColor Yellow
Write-Host "Merge rime\rime_ice.custom.patch.yaml into %APPDATA%\Rime\rime_ice.custom.yaml"
Write-Host "Then right-click Weasel and choose Redeploy."
