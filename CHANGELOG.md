# Changelog

## v2.2.0 - 2026-09-08

- 模型驻留时间调整为 2 分钟。
- debounce 从 500 ms 降至 250 ms。
- Bridge 文件轮询从 200 ms 降至 50 ms。
- 默认限制生成 32 tokens，temperature 0.1。
- 自启动和托盘资源释放保持可用。
- 延长 owned Ollama Server 空闲退出时间以匹配模型驻留策略。

## v2.1.0 - 2026-09-08

- 修复 Lua filter 中同步 `os.execute()` 导致的输入卡死。
- 改为流式 yield 候选，避免收集完整候选列表。

## v2.0.0 - 2026-09-08

- 支持中文 ↔ 英文自动识别。
- 增加开机自动启动与托盘管理。

## v1.0.0 - 2026-09-08

- 将 PowerShell Bridge 重写为无终端窗口的 Windows EXE。
