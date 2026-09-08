# Release guide

## 首次发布建议

1. 在 GitHub 创建一个空仓库，例如 `RimeTranslate`。
2. push 当前源码到 `main`。
3. 检查 Actions 中 `Build` 成功。
4. 创建版本 tag：

```powershell
git tag -a v2.2.0 -m "RimeTranslate v2.2.0"
git push origin v2.2.0
```

5. `Release` workflow 会自动：
   - 在 Windows runner 编译 GUI EXE；
   - 打包 Lua、安装脚本、README、许可证说明；
   - 生成 SHA256；
   - 创建 GitHub Release。

## 仓库中不要提交

- Weasel / rime-ice / Ollama / RIMES 的克隆源码；
- Gemma 模型权重；
- `%APPDATA%\\Rime` 用户词库；
- `dist/` 本地二进制与 zip；
- 日志和本机配置。
