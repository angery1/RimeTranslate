$ErrorActionPreference = "SilentlyContinue"

function Show-State($ok, $label, $detail) {
    if ($ok) { Write-Host "[OK]   $label - $detail" -ForegroundColor Green }
    else { Write-Host "[MISS] $label - $detail" -ForegroundColor Yellow }
}

$rime = Join-Path $env:APPDATA "Rime"
Show-State (Test-Path $rime) "Rime user directory" $rime
Show-State (Test-Path (Join-Path $rime "rime_ice.custom.yaml")) "rime-ice custom config" (Join-Path $rime "rime_ice.custom.yaml")

$ollama = Get-Command ollama -ErrorAction SilentlyContinue
Show-State ($null -ne $ollama) "Ollama CLI" ($(if ($ollama) { $ollama.Source } else { "not in PATH" }))

if ($ollama) {
    $models = & ollama list 2>$null | Out-String
    Show-State ($models -match "gemma3:1b") "gemma3:1b" "ollama model"
}

$luaDir = Join-Path $rime "lua"
Show-State (Test-Path (Join-Path $luaDir "async_ollama_filter.lua")) "translation filter" (Join-Path $luaDir "async_ollama_filter.lua")
Show-State (Test-Path (Join-Path $luaDir "async_refresh.lua")) "refresh processor" (Join-Path $luaDir "async_refresh.lua")
