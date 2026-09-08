param(
    [string]$WorkspaceRoot = "E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法"
)

$ErrorActionPreference = "Stop"
$projectDir = Join-Path $WorkspaceRoot "RimeTranslate"
$depsDir = Join-Path $WorkspaceRoot "第三方开源仓库"
$releaseDir = Join-Path $WorkspaceRoot "发布包"

New-Item -ItemType Directory -Force -Path $WorkspaceRoot, $projectDir, $depsDir, $releaseDir | Out-Null

Write-Host "Workspace created:" -ForegroundColor Green
Write-Host "  本项目:        $projectDir"
Write-Host "  第三方开源仓库: $depsDir"
Write-Host "  发布包:        $releaseDir"
Write-Host ""
Write-Host "Copy/extract this GitHub source package into the 本项目 directory."
Write-Host "Run scripts\clone-references.ps1 only if you want local copies of upstream repositories."
