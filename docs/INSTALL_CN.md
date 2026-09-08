# RimeTranslate 详细安装说明

[返回中文 README](../README.md) | [English installation guide](INSTALL_EN.md)

这份文档用于解决 README 中没有展开的安装细节，尤其是**目录解压**和 **Rime 自定义配置**。

## 1. 安装前提

先安装：

- 小狼毫 Weasel；
- 雾凇拼音 rime-ice；
- Ollama。

然后：

```powershell
ollama pull gemma3:1b
```

## 2. 下载 Release，而不是源码 ZIP

下载：

```text
RimeTranslate-vX.Y.Z-windows-x64.zip
```

不要下载 `Code → Download ZIP` 作为普通安装包。

## 3. 正确解压

Release ZIP 内部已经有：

```text
RimeTranslate\
```

因此推荐直接解压到：

```text
C:\Users\你的用户名\
```

最终：

```text
C:\Users\你的用户名\RimeTranslate\
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

### 常见错误

```text
C:\Users\你\RimeTranslate\RimeTranslate\RimeTranslate.exe
```

出现这种结构说明多套了一层目录。把最里面 `RimeTranslate` 中的所有内容移动到外层即可。

## 4. 双击 Install.cmd

安装器会：

1. 创建 `%APPDATA%\Rime\lua`（如果还不存在）；
2. 备份并安装 `async_ollama_filter.lua`、`async_refresh.lua`；
3. 注册当前用户开机自启动；
4. 检查 `%APPDATA%\Rime\rime_ice.custom.yaml`；
5. 启动 `RimeTranslate.exe`。

## 5. 没有 rime_ice.custom.yaml 怎么办？

**不用手工创建。**

新版安装器检测到该文件不存在时，会直接使用项目提供的配置模板创建：

```text
%APPDATA%\Rime\rime_ice.custom.yaml
```

此时你只需要：

```text
小狼毫 → 重新部署
```

## 已有 rime_ice.custom.yaml 怎么办？

如果你以前已经自定义过雾凇拼音，安装器**不会覆盖**你的文件。

它会保留原配置，并把待合并模板复制到：

```text
%APPDATA%\Rime\RimeTranslate.patch.to_merge.yaml
```

项目需要加入的内容本质上是：

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

### 最容易理解的合并例子

假设你原来有：

```yaml
patch:
  menu/page_size: 9
```

不要写成两个 `patch:`：

```yaml
patch:
  menu/page_size: 9

patch:
  "switches/+":
    ...
```

应该保留**一个** `patch:`：

```yaml
patch:
  menu/page_size: 9

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

如果你的文件里已经存在同名的 `"switches/+"`、`"key_binder/bindings/+"` 等键，不要简单再复制一个同名键；应把对应列表项合并到原键下面。

完成后：

```text
小狼毫 → 重新部署
```

## 6. 测试

按：

```text
Ctrl + Shift + Y
```

测试中文：

```text
nihaoshijie
```

预期：

```text
1 你好世界
2 Hello world 🌐
```

测试英文时保持雾凇拼音**中文模式**，输入：

```text
hello
```

预期：

```text
1 hello
2 你好 🌐
```

## 7. 常见问题

### 开启翻译但没有结果

检查：

```powershell
ollama list
where.exe ollama
```

并查看：

```text
%TEMP%\rime_ollama_bridge\bridge.log
```

### 第一次翻译慢

模型冷启动需要加载 `gemma3:1b`。后续约 2 分钟驻留期间通常会更快。

### 要运行 Vivado

托盘右键：

```text
释放翻译资源
```

即可立即卸载模型。
