# RimeTranslate

[中文](README.md) | **English**

> Show **bidirectional Chinese ↔ English translation directly inside the Rime/Weasel candidate list on Windows**.  
> Chinese candidate → English; English candidate → Chinese. Translation is powered by a local Ollama model through an asynchronous bridge.

[Latest Release](https://github.com/angery1/RimeTranslate/releases/latest) · [Detailed installation](docs/INSTALL_EN.md) · [Model guide](docs/MODELS.md)

## Real screenshots

All screenshots below are from the project running on Windows.

### Chinese → English

![Chinese to English: 基本设置 / Basic settings](docs/images/candidate-zh-to-en.png)

### English → Chinese

Stay in the rime-ice Chinese input mode so the English text becomes a Rime candidate:

![English to Chinese: infrastructure / 基础设施](docs/images/candidate-en-to-zh.png)

### Translation toggle

Use `Ctrl + Shift + Y` to switch translation on/off:

![Translation toggle](docs/images/translation-toggle.png)

## Key features

- Automatic Chinese ↔ English direction detection;
- translated result appears directly in the Rime candidate list;
- asynchronous Ollama bridge, so AI inference does not synchronously block normal Rime candidate generation;
- local inference by default, without relying on an online translation website;
- `gemma3:1b` remains loaded for about 2 minutes to reduce repeated cold starts;
- tray command for immediately releasing translation resources;
- per-user Windows startup support.

## Quick start

### 1. Install prerequisites

Install separately:

1. [Rime Weasel](https://github.com/rime/weasel)
2. [rime-ice](https://github.com/iDvel/rime-ice)
3. [Ollama](https://github.com/ollama/ollama)

Then download the default model:

```powershell
ollama pull gemma3:1b
```

> Ollama must be installed, but you do **not** need to keep an `ollama run` terminal open or enable Ollama at Windows startup. RimeTranslate connects to or starts Ollama when translation is needed.

### 2. Download and extract the Release

Open [Releases](https://github.com/angery1/RimeTranslate/releases/latest) and download:

```text
RimeTranslate-vX.Y.Z-windows-x64.zip
```

Do not use `Code → Download ZIP`; that is the source archive.

The Release ZIP already contains a top-level `RimeTranslate` directory. The easiest method is to extract it directly to:

```text
C:\Users\YOUR_NAME\
```

The final layout should be:

```text
C:\Users\YOUR_NAME\RimeTranslate\
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
   ├─ INSTALL_EN.md
   ├─ MODELS.md
   └─ images\
```

If `RimeTranslate.exe` and `Install.cmd` are immediately visible in `C:\Users\YOUR_NAME\RimeTranslate\`, extraction is correct. Avoid nested layouts such as `RimeTranslate\RimeTranslate\...`.

### 3. Run the installer

Double-click:

```text
Install.cmd
```

The installer copies the Lua files, backs up files with the same names, enables startup for RimeTranslate, and starts the bridge. If `rime_ice.custom.yaml` does not exist, it creates the required configuration automatically. Existing personalized YAML files are never overwritten; see [the detailed installation guide](docs/INSTALL_EN.md).

### 4. Redeploy Weasel

```text
Weasel → Deploy / Redeploy
```

### 5. Use it

Press:

```text
Ctrl + Shift + Y
```

to toggle translation.

Chinese → English: type Pinyin normally.  
English → Chinese: stay in the **rime-ice Chinese input mode** so English is handled as a Rime candidate. Pure ASCII passthrough bypasses the candidate filter and therefore cannot receive a translation candidate.

## Model recommendations

The current Release defaults to **`gemma3:1b`**. Package sizes below come from the Ollama model library. “Suggested free system memory” is an **engineering guideline for this short-text use case, not an official requirement**. Actual RAM/VRAM usage varies with quantization, context length, CPU/GPU offload, and Ollama version.

| Model | Ollama package | Suggested free RAM | Character | Recommendation |
| --- | ---: | ---: | --- | --- |
| `gemma3:270m` | 292 MB | 1–2 GB | lightest / fastest, weaker quality | very low-end systems |
| `qwen2.5:0.5b` | 398 MB | 2–3 GB | lightweight, good fit for Chinese/English | low-resource pick |
| **`gemma3:1b`** | **815 MB** | **3–4 GB** | balanced speed / quality / memory | **default pick** |
| `qwen2.5:1.5b` | 986 MB | 4–6 GB | stronger quality focus | quality-oriented pick |
| `qwen2.5:3b` | 1.9 GB | 6–8 GB | better quality, more latency and memory | resource-rich systems |
| `gemma3:4b` | 3.3 GB | 8–12 GB | heavier; limited benefit for an IME | not a default choice |

CPU-only execution is supported. When a compatible GPU is available, Ollama can offload model work to GPU. Run:

```powershell
ollama ps
```

to inspect whether a loaded model is using `100% CPU`, `100% GPU`, or a CPU/GPU split.

> If you also use Vivado, synthesis, or large simulations, `gemma3:1b` or `qwen2.5:0.5b` is preferable. Before a memory-heavy task, right-click the Rime Translate tray icon and choose **Release translation resources**.

The current binary has the model name in `src/main.go` as `modelName`; using another model currently requires changing that value and rebuilding. See [docs/MODELS.md](docs/MODELS.md).

## How it works

```text
Weasel / Rime / rime-ice
        ↓
async_ollama_filter.lua
        ↓ request.txt
RimeTranslate.exe
        ↓
Ollama + local model
        ↓ response.txt
async_refresh.lua + candidate refresh
        ↓
translated candidate
```

See [Architecture](docs/ARCHITECTURE.md) for details.

## Extensibility

RimeTranslate is not fundamentally limited to translation. Its core pipeline is:

```text
Rime candidate → local background bridge → local AI model → inject result back into the candidate list
```

Because the backend is a local AI model, the project can be extended by changing the prompt, model, and result handling—for example multilingual translation, English polishing, grammar correction, shortening/expanding text, technical-term explanations, fixed-style rewriting, or other Ollama/local inference backends.

**Chinese ↔ English translation is only the first application of this architecture.**

## More documentation

- [Detailed installation](docs/INSTALL_EN.md)
- [Model selection and resource guide](docs/MODELS.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Development](docs/DEVELOPMENT.md)
- [Changelog](CHANGELOG.md)
- [Third-party notices](THIRD_PARTY_NOTICES.md)

## Known limitations

- Currently focused on Windows x64;
- English → Chinese requires the active Rime schema to generate English candidates;
- full English sentences containing spaces are not yet ideal for the current candidate pipeline;
- the first model load is usually slower than translations while the model remains resident;
- Lua ↔ Bridge IPC currently uses local files and may later move to Named Pipes or another local IPC mechanism.

## License

Project-owned code, Lua, scripts, and documentation are licensed under the [MIT License](LICENSE). Weasel, rime-ice, Ollama, Gemma, and other third-party components keep their own licenses.
