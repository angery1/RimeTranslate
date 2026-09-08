# Third-Party Notices / 第三方项目边界

本仓库的原则是：**第三方源码、词库、二进制和模型权重不直接进入本项目 Git 仓库。**

| 项目 | 角色 | 上游 | 许可证/条款 | 是否打包 |
|---|---|---|---|---|
| Rime Weasel / 小狼毫 | Windows Rime 输入法前端；运行时前置条件 | https://github.com/rime/weasel | GPL-3.0 | 否 |
| rime-ice / 雾凇拼音 | 当前推荐 Rime 输入方案与词库 | https://github.com/iDvel/rime-ice | GPL-3.0-only | 否 |
| librime-lua | Lua processor/filter API 上游 | https://github.com/hchunhui/librime-lua | BSD-3-Clause | 否 |
| Ollama | 本地模型服务 / REST API | https://github.com/ollama/ollama | MIT | 否 |
| Gemma 3 (`gemma3:1b`) | 默认翻译模型 | https://ollama.com/library/gemma3 | Gemma Terms of Use | 否；由用户自行 `ollama pull` |
| RIMES | 架构参考（非运行依赖） | https://github.com/scholay/rimes | RIMES 自有代码 MIT；其 bundled Rime data 保留各自许可证 | 否 |

## 重要说明

1. 本项目通过本机 HTTP API 调用 Ollama，不包含 Ollama 源码或二进制。
2. 本项目不包含 Gemma 模型权重；用户通过 Ollama 自行下载，模型受其自身条款约束。
3. 本项目不复制雾凇拼音词库或 Weasel 源码，仅在用户已安装的 Rime 环境中安装本项目自有 Lua 脚本并要求用户合并配置 patch。
4. `scholay/rimes` 仅用于架构学习与设计参考，本项目没有复制其源码。
5. 本文件仅用于说明项目边界，不构成法律意见。如你未来开始直接打包第三方二进制、词库或源代码，需要重新核对对应许可证的再分发义务。
