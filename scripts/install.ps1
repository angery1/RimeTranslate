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
$patchTemplate = Join-Path $root "rime\rime_ice.custom.patch.yaml"
$customYaml = Join-Path $rime "rime_ice.custom.yaml"
$mergeHint = Join-Path $rime "RimeTranslate.patch.to_merge.yaml"

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

$configState = "manual"
if (-not (Test-Path $customYaml)) {
    Copy-Item $patchTemplate $customYaml -Force
    $configState = "created"
} else {
    Copy-Item $customYaml (Join-Path $backup "rime_ice.custom.yaml") -Force
    $raw = Get-Content -Path $customYaml -Raw -Encoding UTF8
    $hasSwitch = $raw -match "ollama_translation"
    $hasProcessor = $raw -match "async_refresh"
    $hasFilter = $raw -match "async_ollama_filter"

    if ($hasSwitch -and $hasProcessor -and $hasFilter) {
        $configState = "ready"
        Remove-Item $mergeHint -Force -ErrorAction SilentlyContinue
    } else {
        Copy-Item $patchTemplate $mergeHint -Force
        $configState = "manual"
    }
}

$runKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
New-Item -Path $runKey -Force | Out-Null
Set-ItemProperty -Path $runKey -Name "RimeTranslate" -Value ('"' + $exe + '"')

Get-Process -Name "RimeTranslate" -ErrorAction SilentlyContinue | Stop-Process -Force -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 300
Start-Process -FilePath $exe

Write-Host ""
Write-Host "RimeTranslate installed." -ForegroundColor Green
Write-Host "Lua/config backup: $backup"
Write-Host "Startup: enabled for current user"
Write-Host ""

switch ($configState) {
    "created" {
        Write-Host "Rime config: created automatically." -ForegroundColor Green
        Write-Host "NEXT STEP: right-click Weasel and choose Redeploy." -ForegroundColor Yellow
    }
    "ready" {
        Write-Host "Rime config: RimeTranslate entries already detected." -ForegroundColor Green
        Write-Host "NEXT STEP: right-click Weasel and choose Redeploy." -ForegroundColor Yellow
    }
    default {
        Write-Host "Rime config: an existing custom file was detected and was NOT overwritten." -ForegroundColor Yellow
        Write-Host "Merge the template below into your existing rime_ice.custom.yaml:" -ForegroundColor Yellow
        Write-Host "  $mergeHint"
        Write-Host "Detailed guide: docs\INSTALL_CN.md / docs\INSTALL_EN.md"
        Write-Host "Then right-click Weasel and choose Redeploy."
    }
}
