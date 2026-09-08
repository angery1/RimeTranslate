# RimeTranslate · Windows 双语输入法候选翻译

RimeTranslate 是一个面向 Windows 小狼毫（Weasel）+ Rime 的轻量本地翻译桥接器。

目标体验：

```text
nihaoshijie

1 你好世界
2 Hello world 🌐
```

以及在雾凇拼音中文模式中输入英文候选时：

```text
hello

1 hello
2 你好 🌐
```

它不会把 LLM 同步塞进 Rime 的候选线程，而是采用异步桥接：Rime 先正常给出候选，后台程序调用 Ollama，翻译完成后再触发候选刷新。因此即使模型较慢，正常中文输入也不应被模型推理阻塞。

## 当前版本

**v2.2.0**

默认参数：

- 模型：`gemma3:1b`
- 中文 ↔ 英文自动方向判断
- 防抖：250 ms
- Bridge 文件轮询：50 ms
- 模型 `keep_alive`：2 分钟
- 最大生成：32 tokens
- temperature：0.1
- 若 Ollama 由本程序按需启动：约 3 分钟空闲后关闭该服务
- 托盘菜单支持“释放翻译资源”
- 支持当前用户开机自启动

## 架构

```text
键盘
  ↓
小狼毫 / Rime / 雾凇拼音
  ↓
async_ollama_filter.lua
  ↓  request.txt
RimeTranslate.exe
  ↓  HTTP 127.0.0.1:11434
Ollama + gemma3:1b
  ↓  response.txt
async_ollama_filter.lua
  ↓
英文/中文翻译作为候选 #2
```

翻译完成后，Bridge 发送 F24；`async_refresh.lua` 捕获 F24 并调用 Rime 的 composition refresh，使翻译候选异步出现。

## 依赖

本仓库 **不打包、不复制、不重新发布** 第三方项目源码或模型权重。用户自行安装：

1. [Rime Weasel / 小狼毫](https://github.com/rime/weasel)
2. [雾凇拼音 rime-ice](https://github.com/iDvel/rime-ice)（当前推荐方案）
3. [Ollama](https://github.com/ollama/ollama)
4. Ollama 模型：`gemma3:1b`

Lua 能力依赖 Rime 环境中的 librime-lua；相关上游见 [hchunhui/librime-lua](https://github.com/hchunhui/librime-lua)。

RIMES (`scholay/rimes`) 仅作为原生输入法/单进程服务化架构的**设计参考**，不是本项目运行依赖。

详细许可证边界见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。

## 安装

### 1. 安装运行环境

安装小狼毫和雾凇拼音，然后安装 Ollama：

```powershell
ollama pull gemma3:1b
```

无需保持 `ollama run` 终端常开。RimeTranslate 在收到真正的翻译请求时会检查 `127.0.0.1:11434`；若 Ollama 未运行，会尝试后台执行 `ollama serve`。

### 2. 安装 Rime Lua 文件

从 GitHub Release 下载 Windows x64 ZIP，解压到固定目录，例如：

```text
C:\Users\你的用户名\RimeTranslate\
```

运行：

```text
Install.cmd
```

安装器会：

- 复制 `async_ollama_filter.lua` 和 `async_refresh.lua` 到 `%APPDATA%\Rime\lua`
- 备份已有同名 Lua
- 将 `RimeTranslate.exe` 注册到当前用户开机启动项
- 启动 Bridge

### 3. 合并 Rime patch

打开：

```text
%APPDATA%\Rime\rime_ice.custom.yaml
```

将 [`rime/rime_ice.custom.patch.yaml`](rime/rime_ice.custom.patch.yaml) 中 `patch:` 下对应键合并进去。

> 不建议直接覆盖你现有的 `rime_ice.custom.yaml`，因为该文件可能包含你的个人配置。

完成后：

```text
小狼毫 → 重新部署
```

### 4. 使用

`Ctrl + Shift + Y`：切换“翻译关 / 翻译开”。

中文 → 英文：正常输入拼音即可。

英文 → 中文：保持在**雾凇拼音中文模式**，让英文先成为 Rime 候选；如果切到纯 ASCII/西文直通模式，文字会绕过候选 Filter，因此不会产生翻译候选。

## 资源占用策略

模型保持 2 分钟主要用于减少连续翻译时的冷启动。

准备运行 Vivado / 大型仿真时，可右键托盘图标：

```text
释放翻译资源
```

程序会立即请求卸载模型；如果 Ollama Server 是由 RimeTranslate 自己启动的，也会结束它。Bridge 自身仍保持后台等待。

## 从源码构建

要求 Go 1.23+：

```powershell
.\scripts\build.ps1
```

输出：

```text
dist\RimeTranslate.exe
```

也可以手动：

```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -trimpath -ldflags="-H=windowsgui -s -w" -o dist\RimeTranslate.exe .\src
```

## 发布

仓库提供 GitHub Actions。推送类似：

```text
v2.2.0
v2.3.0
```

的 tag 后，Release 工作流会在 Windows runner 编译 EXE、组装安装包并创建 GitHub Release。

## 已知限制

- 当前主要针对 Windows x64。
- 英文 → 中文依赖当前 Rime 方案能够生成英文候选；雾凇拼音适合单词/中英混输，但“包含空格的完整英文长句”并不是当前候选链最理想的输入形式。
- 首次调用模型（冷启动）通常比两分钟驻留期间的后续翻译慢。
- 本项目使用本地文件作为 Rime 与 Bridge 的异步 IPC。后续版本可改为 Named Pipe / 本地 IPC。

## License

本仓库**自行编写的代码与脚本**使用 [MIT License](LICENSE)。第三方软件、词库、模型和参考项目保持各自许可证，本仓库不对其重新授权。
