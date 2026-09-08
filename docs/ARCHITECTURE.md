# Architecture

## 目标

RimeTranslate 的核心设计约束：**LLM 慢不能让正常输入慢。**

同步方案如果在 Rime Lua filter 中执行 `curl` / `io.popen()` 并等待 Ollama 返回，会直接阻塞候选生成。当前方案把模型调用移到独立后台 Bridge。

## 数据流

```text
Rime first candidate
      │
      ├─ Lua 写 request.txt（极短操作）
      │
      └─ Lua 立即 yield 原候选

RimeTranslate.exe
      │
      ├─ 250 ms debounce
      ├─ 自动识别 Han / Latin
      ├─ 确认或启动 Ollama
      ├─ POST /api/generate
      ├─ 写 response.tmp → atomic rename → response.txt
      └─ 模拟 F24

async_refresh.lua
      │
      └─ refresh_non_confirmed_composition()

Lua filter 再次运行
      │
      └─ 如果 response source == 当前首候选，插入 Candidate #2
```

## 资源策略

- `gemma3:1b` keep_alive = 2m
- Ollama 若由 Bridge 启动，Bridge 空闲约 3m 后结束该 server
- 托盘“释放翻译资源”可立即卸载模型/结束 owned Ollama server
- Bridge 空闲时不执行模型推理

## 与 RIMES 的关系

RIMES 强调把输入核心、UI 和服务生命周期做成更一体化的原生输入法架构。本项目当前不重写 Windows TSF，而是保留 Weasel，把翻译服务收敛为一个无终端窗口的原生后台 EXE。它属于“渐进式”路线：

1. 当前：文件 IPC + 独立 Bridge
2. 下一步：Named Pipe / event-based IPC
3. 长期：如有必要，再研究完整 TSF 前端

RIMES 仅作为设计参考，不是构建或运行依赖。
