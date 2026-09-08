# RimeTranslate detailed installation guide

[Chinese README](../README.md) | [中文安装说明](INSTALL_CN.md)

This guide expands the two parts that most often confuse first-time users: **where to extract the Release ZIP** and **how the Rime custom YAML is handled**.

## 1. Prerequisites

Install Weasel, rime-ice, and Ollama, then run:

```powershell
ollama pull gemma3:1b
```

## 2. Download a Release package

Download:

```text
RimeTranslate-vX.Y.Z-windows-x64.zip
```

Do not use `Code → Download ZIP` as the normal Windows installation package.

## 3. Extract correctly

The Release archive already contains a top-level `RimeTranslate` folder. Extract the ZIP directly to:

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
├─ LICENSE
├─ THIRD_PARTY_NOTICES.md
├─ docs\
│  ├─ INSTALL_CN.md
│  └─ INSTALL_EN.md
└─ rime\
   ├─ rime_ice.custom.patch.yaml
   └─ lua\
      ├─ async_ollama_filter.lua
      └─ async_refresh.lua
```

Avoid nested layouts such as:

```text
C:\Users\YOU\RimeTranslate\RimeTranslate\RimeTranslate.exe
```

## 4. Run Install.cmd

The installer will:

1. create `%APPDATA%\Rime\lua` if necessary;
2. back up and install the two Lua files;
3. enable current-user startup;
4. inspect `%APPDATA%\Rime\rime_ice.custom.yaml`;
5. start `RimeTranslate.exe`.

## 5. What if rime_ice.custom.yaml does not exist?

You do not need to create it manually. The installer automatically creates:

```text
%APPDATA%\Rime\rime_ice.custom.yaml
```

from the bundled template.

Then redeploy Weasel.

## What if rime_ice.custom.yaml already exists?

The installer will **not overwrite your custom file**. It keeps your file and copies the RimeTranslate template to:

```text
%APPDATA%\Rime\RimeTranslate.patch.to_merge.yaml
```

RimeTranslate needs these entries under the existing `patch:` mapping:

```yaml
patch:
  "switches/+":
    - name: ollama_translation
      reset: 0
      states: [ 翻译关, 翻译开 ]

  "key_binder/bindings/+":
    - { when: always, accept: Control+Shift+Y, toggle: ollama_translation }

  "engine/processors/@before 0":
    lua_processor@*async_refresh

  "engine/filters/@before last":
    lua_filter@*async_ollama_filter
```

Keep a single top-level `patch:` mapping. If your existing YAML already contains one of the same path keys, merge the list/value into that key rather than creating a duplicate YAML key.

After editing, redeploy Weasel.

## 6. Test

Press `Ctrl + Shift + Y`, then try:

```text
nihaoshijie
```

Expected:

```text
1 你好世界
2 Hello world 🌐
```

For English → Chinese, stay in rime-ice Chinese mode and type:

```text
hello
```

Expected:

```text
1 hello
2 你好 🌐
```

## 7. Troubleshooting

Check Ollama:

```powershell
ollama list
where.exe ollama
```

Bridge log:

```text
%TEMP%\rime_ollama_bridge\bridge.log
```

The first translation after a cold start may be slower. The model normally stays loaded for about 2 minutes to speed up continued use.
