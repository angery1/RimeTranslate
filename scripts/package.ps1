param([string]$Version = "2.2.1")

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
& (Join-Path $root "scripts\build.ps1")

$stage = Join-Path $root "dist\package"
Remove-Item $stage -Recurse -Force -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path (Join-Path $stage "RimeTranslate") | Out-Null
$pkg = Join-Path $stage "RimeTranslate"

Copy-Item (Join-Path $root "dist\RimeTranslate.exe") $pkg
Copy-Item (Join-Path $root "rime") $pkg -Recurse
Copy-Item (Join-Path $root "scripts\install.ps1") $pkg
Copy-Item (Join-Path $root "scripts\check-environment.ps1") $pkg
Copy-Item (Join-Path $root "README.md") $pkg
Copy-Item (Join-Path $root "README_EN.md") $pkg
Copy-Item (Join-Path $root "LICENSE") $pkg
Copy-Item (Join-Path $root "THIRD_PARTY_NOTICES.md") $pkg

New-Item -ItemType Directory -Force -Path (Join-Path $pkg "docs") | Out-Null
Copy-Item (Join-Path $root "docs\INSTALL_CN.md") (Join-Path $pkg "docs")
Copy-Item (Join-Path $root "docs\INSTALL_EN.md") (Join-Path $pkg "docs")
Copy-Item (Join-Path $root "docs\MODELS.md") (Join-Path $pkg "docs")
New-Item -ItemType Directory -Force -Path (Join-Path $pkg "docs\images") | Out-Null
Copy-Item (Join-Path $root "docs\images\*") (Join-Path $pkg "docs\images") -Recurse

$installCmd = @'
@echo off
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0install.ps1"
pause
'@
$installCmd | Set-Content -Path (Join-Path $pkg "Install.cmd") -Encoding ASCII

$zip = Join-Path $root ("dist\RimeTranslate-v" + $Version + "-windows-x64.zip")
Remove-Item $zip -Force -ErrorAction SilentlyContinue
Compress-Archive -Path $pkg -DestinationPath $zip -CompressionLevel Optimal
Write-Host "Package: $zip" -ForegroundColor Green
