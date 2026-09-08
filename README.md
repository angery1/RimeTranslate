# RimeTranslate

**中文** | [English](README_EN.md)

> 在 Windows 小狼毫 / Rime 的**候选框里直接显示中英双向翻译**。  
> 中文候选 → 英文；英文候选 → 中文。翻译由本地 Ollama 模型完成，正常输入与 AI 推理解耦，避免把模型同步阻塞塞进 Rime 候选线程。

[下载最新 Release](https://github.com/angery1/RimeTranslate/releases/latest) · [详细安装说明](docs/INSTALL_CN.md) · [第三方依赖与许可证](THIRD_PARTY_NOTICES.md)

## 效果预览

| 中文 → 英文 | 英文 → 中文 |
| --- | --- |
| `nihaoshijie` → `1 你好世界` → `2 Hello world 🌐` | `hello` → `1 hello` → `2 你好 🌐` |

> README 将使用**真实 Windows 功能截图**展示效果。截图文件建议放在 `docs/images/`；为避免用伪造界面冒充真实运行效果，当前先保留上面的文字预览。截图补齐后可直接替换为图片。

## 你需要先准备

本项目只提供桥接层，不把第三方程序或模型权重打包进仓库。请先安装：

1. [小狼毫 Weasel](https://github.com/rime/weasel)
2. [雾凇拼音 rime-ice](https://github.com/iDvel/rime-ice)
3. [Ollama](https://github.com/ollama/ollama)
4. 模型 `gemma3:1b`

安装 Ollama 后执行：

```powershell
ollama pull gemma3:1b
```

不需要长期打开 `ollama run gemma3:1b`。RimeTranslate 会在需要翻译时按需连接或启动 Ollama。

## 快速安装

### 1. 下载发布包

到 [Releases](https://github.com/angery1/RimeTranslate/releases/latest) 下载：

```text
RimeTranslate-vX.Y.Z-windows-x64.zip
```

> 不要使用 `Code → Download ZIP`，那是源码包，不是 Windows 安装包。

### 2. 解压到正确位置

Release ZIP **内部已经有一个 `RimeTranslate` 顶层文件夹**。

推荐把 ZIP 直接解压到：

```text
C:\Users\你的用户名\
```

最终必须得到：

```text
C:\Users\你的用户名\RimeTranslate\
├─ RimeTranslate.exe
├─ Install.cmd
├─ install.ps1
├─ check-environment.ps1
├─ README.md
├─ README_EN.md
├─ rime\
│  ├─ rime_ice.custom.patch.yaml
│  └─ lua\
│     ├─ async_ollama_filter.lua
│     └─ async_refresh.lua
└─ docs\
   ├─ INSTALL_CN.md
   └─ INSTALL_EN.md
```

**判断是否解压正确：**打开 `C:\Users\你的用户名\RimeTranslate\` 后，应当立刻看到 `RimeTranslate.exe` 和 `Install.cmd`。

不要出现：

```text
RimeTranslate\RimeTranslate\RimeTranslate.exe
```

### 3. 双击安装

双击：

```text
Install.cmd
```

安装器会自动：

- 安装两个 Rime Lua 文件；
- 备份已有同名 Lua；
- 注册 `RimeTranslate.exe` 为当前用户开机自启动；
- 启动后台 Bridge；
- **如果你没有 `rime_ice.custom.yaml`，自动创建所需配置。**

如果你已经有自己的 `rime_ice.custom.yaml`，安装器不会覆盖它，而会提示你按 [详细安装说明](docs/INSTALL_CN.md#已有-rime_icecustomyaml-怎么办) 合并配置。

### 4. 重新部署小狼毫

```text
小狼毫 → 重新部署
```

### 5. 开始使用

按：

```text
Ctrl + Shift + Y
```

切换：

```text
翻译关 ↔ 翻译开
```

中文 → 英文：正常输入拼音即可。

英文 → 中文：保持在**雾凇拼音中文模式**，让英文先进入 Rime 候选；纯 ASCII/西文直通模式会绕过候选 Filter，因此不会产生翻译候选。

## 主要功能

- 中文 ↔ 英文自动判断方向；
- 翻译结果插入候选框，不弹额外翻译窗口；
- 异步调用 Ollama，模型慢不会要求 Rime 同步等待；
- `gemma3:1b` 默认驻留 2 分钟，减少连续翻译冷启动；
- 托盘可立即“释放翻译资源”；
- 支持当前 Windows 用户开机自启动；
- 默认本地推理，不依赖在线翻译网站。

## 资源占用

连续使用时模型默认保持约 **2 分钟**，减少反复加载造成的等待。

运行 Vivado、仿真或其他大型任务前，可以：

```text
右键 Rime Translate 托盘图标
→ 释放翻译资源
```

模型会被立即卸载；如果 Ollama Server 是由 RimeTranslate 自己启动的，也会按程序策略释放。

## 工作原理

```text
小狼毫 / Rime / 雾凇拼音
        ↓
async_ollama_filter.lua
        ↓ request.txt
RimeTranslate.exe
        ↓
Ollama + gemma3:1b
        ↓ response.txt
async_refresh.lua + 候选刷新
        ↓
翻译结果作为候选出现
```

更详细的实现说明见 [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)。

## 可扩展性

RimeTranslate 的核心并不是“固定写死的翻译器”，而是一个：

```text
Rime 候选
→ 本地后台桥接器
→ 本地 AI 模型
→ 结果重新注入候选框
```

的通用框架。

因为后端使用的是**本地 AI 模型**，后续可以通过修改 Prompt、模型和结果处理逻辑，把它扩展成更多候选增强功能，例如：

- 中英以外的多语言翻译；
- 英文润色、语法改写；
- 输入内容的简写 / 扩写；
- 专业术语解释；
- 中英文专业词汇转换；
- 自定义固定风格改写；
- 接入其他 Ollama 模型或本地推理后端。

也就是说，当前“中英双向翻译”只是这个架构的第一个应用场景。

## 更多文档

- [详细安装说明](docs/INSTALL_CN.md)
- [架构说明](docs/ARCHITECTURE.md)
- [开发说明](docs/DEVELOPMENT.md)
- [更新记录](CHANGELOG.md)
- [第三方依赖与许可证](THIRD_PARTY_NOTICES.md)

## 已知限制

- 当前主要面向 Windows x64；
- 英文 → 中文依赖当前 Rime 方案能够产生英文候选；
- 包含空格的完整英文长句还不是当前候选链最理想的输入形式；
- 第一次加载模型通常比模型驻留期间的后续翻译慢；
- 当前 Lua ↔ Bridge 使用本地文件异步通信，后续可进一步改为 Named Pipe / 本地 IPC。

## License

本仓库自行编写的代码、Lua、脚本与文档使用 [MIT License](LICENSE)。Weasel、rime-ice、Ollama、Gemma 等第三方项目继续遵循各自许可证，本仓库不对其重新授权。
