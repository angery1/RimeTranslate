# Development

## 本地开发要求

- Windows 10/11 x64
- Go 1.23+
- Rime Weasel
- rime-ice（推荐测试方案）
- Ollama
- `gemma3:1b`

## 开发构建

```powershell
.\scripts\build.ps1
```

## 手动测试矩阵

1. 翻译关闭：拼音输入无额外延迟。
2. 翻译开启：中文首候选立即出现，翻译稍后出现。
3. `nihao` → `你好` + `Hello`。
4. 中文模式中 `hello` → `hello` + `你好`（取决于英文候选）。
5. Ollama 未启动时：Bridge 按需启动。
6. Ollama 已由用户启动时：Bridge 不应在释放资源时杀掉外部服务。
7. 托盘“释放翻译资源”：模型卸载。
8. 退出 Bridge 后：Rime 正常输入仍能工作，只是没有翻译候选。
9. Vivado 等高负载应用同时运行时：翻译关闭不应触发模型。

## 日志

```text
%TEMP%\rime_ollama_bridge\bridge.log
```
