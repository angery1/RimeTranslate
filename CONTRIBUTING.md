# Contributing

欢迎提交 issue / pull request。

基本要求：

- 不要把 Weasel、rime-ice、Ollama、Gemma 模型权重直接提交进本仓库。
- 不要提交 `%APPDATA%\Rime` 中的个人词库、用户数据或隐私文件。
- 修改 Lua 时优先保证：候选线程不执行同步外部进程或长耗时网络/模型调用。
- Windows Bridge 修改后至少执行 `scripts\build.ps1` 和 README 中的手动测试矩阵。
- 新功能请同时更新 CHANGELOG / 文档。
