# Changelog

## v2.2.1 - 2026-09-10

- 重构 README，将功能、下载、快速安装放到前面，详细说明移入 docs。
- README 增加中文 / English 双语切换。
- README 增加项目在 Windows 上的真实运行截图：中文 → 英文、英文 → 中文、翻译开关。
- 新增 `docs/MODELS.md`，说明不同本地模型的模型包大小、资源取舍和部署建议。
- 默认推荐 `gemma3:1b`；低资源场景推荐 `qwen2.5:0.5b`；质量优先可考虑 `qwen2.5:1.5b`。
- Release 包同时携带 README.md、README_EN.md、中英文安装文档、模型说明和功能截图。
- 首次安装且不存在 `rime_ice.custom.yaml` 时自动创建 RimeTranslate 配置。
- 已存在个人 `rime_ice.custom.yaml` 时不覆盖，生成 `RimeTranslate.patch.to_merge.yaml` 供手工合并。
- 增加本地 AI 架构可扩展性说明。

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
