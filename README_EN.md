# RimeTranslate

[中文](README.md) | **English**

> Show **bidirectional Chinese ↔ English translations directly in the Rime/Weasel candidate list on Windows**.  
> Chinese candidate → English; English candidate → Chinese. Translation runs through a local Ollama model and is kept off the synchronous Rime candidate-generation path.

[Latest Release](https://github.com/angery1/RimeTranslate/releases/latest) · [Detailed installation](docs/INSTALL_EN.md) · [Third-party notices](THIRD_PARTY_NOTICES.md)

## Preview

| Chinese → English | English → Chinese |
| --- | --- |
| `nihaoshijie` → `1 你好世界` → `2 Hello world 🌐` | `hello` → `1 hello` → `2 你好 🌐` |

> The README is intended to use **real Windows screenshots**. To avoid presenting a fabricated UI as a real run, the repository currently keeps the text preview above until real screenshots are added under `docs/images/`.

## Prerequisites

Install these separately:

1. [Rime Weasel](https://github.com/rime/weasel)
2. [rime-ice](https://github.com/iDvel/rime-ice)
3. [Ollama](https://github.com/ollama/ollama)
4. `gemma3:1b`

Then run:

```powershell
ollama pull gemma3:1b
```

You do not need to keep an `ollama run` terminal open. RimeTranslate connects to or starts Ollama on demand.

## Quick start

### 1. Download the Release package

Open [Releases](https://github.com/angery1/RimeTranslate/releases/latest) and download:

```text
RimeTranslate-vX.Y.Z-windows-x64.zip
```

Do **not** use `Code → Download ZIP`; that is the source archive.

### 2. Extract to the correct directory

The Release ZIP already contains a top-level `RimeTranslate` folder. The easiest method is to extract the ZIP directly to:

```text
C:\Users\YOUR_NAME\
```

The result must look like this:

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
   └─ INSTALL_EN.md
```

If `RimeTranslate.exe` and `Install.cmd` are immediately visible inside that directory, extraction is correct.

Avoid nested layouts such as:

```text
RimeTranslate\RimeTranslate\RimeTranslate.exe
```

### 3. Run the installer

Double-click:

```text
Install.cmd
```

It will:

- install the two Rime Lua files;
- back up existing files with the same names;
- enable current-user startup for `RimeTranslate.exe`;
- start the background bridge;
- **create the required `rime_ice.custom.yaml` automatically if you do not already have one.**

If you already have a customized `rime_ice.custom.yaml`, the installer will not overwrite it. Follow [the detailed installation guide](docs/INSTALL_EN.md#what-if-rime_icecustomyaml-already-exists) instead.

### 4. Redeploy Weasel

```text
Weasel → Deploy / Redeploy
```

### 5. Use it

Press:

```text
Ctrl + Shift + Y
```

to toggle translation on/off.

Chinese → English: type Pinyin normally.

English → Chinese: stay in the **rime-ice Chinese input mode** so English becomes a Rime candidate. Pure ASCII passthrough bypasses the candidate filter and therefore cannot receive a translation candidate.

## Features

- Automatic Chinese ↔ English direction detection;
- translation appears inside the candidate list instead of a separate popup;
- asynchronous Ollama bridge so AI inference does not synchronously block normal candidate generation;
- `gemma3:1b` stays loaded for about 2 minutes by default to reduce repeated cold starts;
- tray command to release translation resources immediately;
- per-user Windows startup support;
- local inference by default.

## Resource usage

The model stays loaded for about **2 minutes** during active use to reduce repeated cold starts.

Before running Vivado, simulations, or other memory-heavy workloads, use:

```text
Right-click Rime Translate tray icon
→ Release translation resources
```

## How it works

```text
Weasel / Rime / rime-ice
        ↓
async_ollama_filter.lua
        ↓ request.txt
RimeTranslate.exe
        ↓
Ollama + gemma3:1b
        ↓ response.txt
async_refresh.lua + candidate refresh
        ↓
translated candidate
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for details.

## Extensibility

RimeTranslate is not fundamentally limited to translation. Its core pipeline is:

```text
Rime candidate
→ local background bridge
→ local AI model
→ inject processed result back into the candidate list
```

Because the backend is a **local AI model**, the project can be modified or extended by changing the prompt, model, and result handling. Possible directions include:

- multilingual translation;
- English polishing and grammar correction;
- shortening or expanding text;
- technical-term explanations;
- domain-specific terminology conversion;
- style rewriting;
- other Ollama models or local inference backends.

Chinese ↔ English translation is therefore only the first application of the architecture.

## More documentation

- [Detailed installation](docs/INSTALL_EN.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Development](docs/DEVELOPMENT.md)
- [Changelog](CHANGELOG.md)
- [Third-party notices](THIRD_PARTY_NOTICES.md)

## Known limitations

- Currently focused on Windows x64;
- English → Chinese requires the active Rime schema to produce English candidates;
- full English sentences containing spaces are not yet ideal for the current candidate pipeline;
- the first model load is usually slower than translations while the model remains resident;
- Lua ↔ Bridge IPC currently uses local files and may later move to Named Pipes or another local IPC mechanism.

## License

Project-owned code, Lua, scripts, and documentation are licensed under the [MIT License](LICENSE). Weasel, rime-ice, Ollama, Gemma, and other third-party components keep their own licenses.
