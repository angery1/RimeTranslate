# 模型选择与硬件资源 / Model & Hardware Guide

RimeTranslate 当前发布版默认使用 `gemma3:1b`。模型由 Ollama 单独安装，仓库和 Release **不包含模型权重**。

## 先理解：模型文件大小 ≠ 实际 RAM / VRAM 占用

Ollama 模型页显示的是下载/存储大小。运行时还需要模型工作内存、KV cache、Ollama 运行时以及 CPU/GPU offload 相关开销，因此实际内存或显存通常会高于模型文件本身。

下面“建议空闲系统内存”是为了帮助 RimeTranslate 用户选型的**工程经验值**，不是 Ollama 或模型厂商给出的硬性最低配置。输入法翻译通常是短文本，实际占用还会受量化版本、上下文长度、操作系统、Ollama 版本和 GPU offload 方式影响。

## 常用模型对比

| 模型 | Ollama 官方模型大小 | 建议空闲系统内存* | 响应倾向 | 本项目建议 |
| --- | ---: | ---: | --- | --- |
| `gemma3:270m` | 292 MB | 1–2 GB | 最快、最轻，翻译能力有限 | 极低配置 / 速度优先 |
| `qwen2.5:0.5b` | 398 MB | 2–3 GB | 很快，对中文和英文较友好 | **低配置推荐** |
| `gemma3:1b` | 815 MB | 3–4 GB | 速度、质量、内存较均衡 | **默认推荐** |
| `qwen2.5:1.5b` | 986 MB | 4–6 GB | 稍慢，中文/英文短文本质量更稳 | **质量优先推荐** |
| `qwen2.5:3b` | 1.9 GB | 6–8 GB | 质量继续提升，但冷启动和推理变慢 | 资源充足时可选 |
| `gemma3:4b` | 3.3 GB | 8–12 GB | 更重，普通输入法翻译收益有限 | 不建议作为默认模型 |
| `qwen2.5:7b` | 4.7 GB | 12–16 GB | 更高质量但明显更重 | 一般不推荐输入法常驻 |

\* 经验建议，不是官方最低要求。

模型大小来源：

- Gemma 3: https://ollama.com/library/gemma3
- Qwen2.5: https://ollama.com/library/qwen2.5

## 推荐怎么选

### 1. 默认：`gemma3:1b`

适合大多数电脑，也是当前官方 Release 的默认配置。

优势：

- 815 MB 模型包，明显小于 3B/4B/7B 级模型；
- 对候选框里的短词、短句翻译速度较合适；
- 当前项目已经围绕它设置了 2 分钟 `keep_alive`；
- 对“输入法 + 日常办公 + 偶尔运行大型软件”的组合比较均衡。

### 2. 内存紧张 / 同时运行 Vivado：`qwen2.5:0.5b`

如果电脑需要把更多内存留给 Vivado、综合或仿真，可以优先尝试这个模型。模型更小，冷启动通常更轻，但复杂句子的质量会下降。

### 3. 更看重翻译质量：`qwen2.5:1.5b`

如果内存更充足，并且希望中英短句翻译更稳，可以尝试 1.5B。它仍远轻于 7B 级模型，但资源占用和首次加载时间会高于默认 1B。

### 4. 为什么通常不推荐 7B+

RimeTranslate 的核心使用场景是输入过程中的**短文本低延迟候选**，而不是长文推理。7B 及以上模型可能提高复杂文本质量，但同时会明显增加：

- 模型加载时间；
- RAM / VRAM 占用；
- CPU-only 推理延迟；
- 与 Vivado、IDE、浏览器等程序争抢资源的概率。

因此对输入法而言，更大的模型不一定意味着更好的综合体验。

## CPU 和 GPU

RimeTranslate 本身不要求独立 GPU。Ollama 可以使用 CPU-only；若系统有兼容 GPU，Ollama 可自动把部分或全部模型计算 offload 到 GPU。

查看当前模型实际运行位置：

```powershell
ollama ps
```

`PROCESSOR` 一栏可能显示：

```text
100% CPU
100% GPU
48%/52% CPU/GPU
```

因此不要仅根据模型下载大小推断显存占用，最好在自己的电脑上通过 `ollama ps` 和任务管理器观察。

## 与 Vivado 等大型软件共用电脑

本项目默认让模型在最后一次翻译后保持约 **2 分钟**，目的是避免每输入几次就重新加载模型。

需要把资源让给 Vivado 时无需等待两分钟：

```text
右键 Rime Translate 托盘图标
→ 释放翻译资源
```

Bridge 会请求立即卸载模型；如果 Ollama Server 是由 RimeTranslate 自己启动的，也会按程序的资源管理策略关闭。

## 安装不同模型

例如：

```powershell
ollama pull qwen2.5:0.5b
ollama pull gemma3:1b
ollama pull qwen2.5:1.5b
```

当前版本的模型名在：

```text
src/main.go
```

中设置：

```go
modelName = "gemma3:1b"
```

要使用其他模型，目前需要改为例如：

```go
modelName = "qwen2.5:0.5b"
```

然后重新构建：

```powershell
.\scripts\build.ps1
```

未来可以把模型选择做成配置文件或托盘设置，避免重新编译。

---

## English summary

The Release currently defaults to **`gemma3:1b`**, which is the recommended balance for short IME translations. For memory-constrained systems, especially machines also running Vivado, consider **`qwen2.5:0.5b`**. For higher translation quality with more available memory, **`qwen2.5:1.5b`** is a reasonable step up.

The package size shown by Ollama is **not** the same as runtime RAM/VRAM usage. Runtime memory also depends on quantization, context/KV cache, Ollama version, and CPU/GPU offload. Use `ollama ps` plus the Windows Task Manager to inspect your actual system.
