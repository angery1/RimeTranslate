param(
    [string]$WorkspaceRoot = "E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法"
)

$ErrorActionPreference = "Stop"
$deps = Join-Path $WorkspaceRoot "第三方开源仓库"
New-Item -ItemType Directory -Force -Path $deps | Out-Null

$repos = @(
    @{ Name = "weasel";       Url = "https://github.com/rime/weasel.git"; Role = "运行时上游/研究" },
    @{ Name = "rime-ice";     Url = "https://github.com/iDvel/rime-ice.git"; Role = "输入方案/词库" },
    @{ Name = "librime-lua";  Url = "https://github.com/hchunhui/librime-lua.git"; Role = "Lua API 参考" },
    @{ Name = "ollama";       Url = "https://github.com/ollama/ollama.git"; Role = "本地模型服务" },
    @{ Name = "rimes";        Url = "https://github.com/scholay/rimes.git"; Role = "架构参考；非依赖" }
)

foreach ($repo in $repos) {
    $target = Join-Path $deps $repo.Name
    if (Test-Path $target) {
        Write-Host "[SKIP] $($repo.Name) already exists - $($repo.Role)"
        continue
    }
    Write-Host "[CLONE] $($repo.Name) - $($repo.Role)"
    git clone --depth 1 $repo.Url $target
}

Write-Host ""
Write-Host "Third-party repositories are stored outside the RimeTranslate Git repository:" -ForegroundColor Green
Write-Host $deps
