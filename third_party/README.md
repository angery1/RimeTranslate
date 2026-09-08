# Third-party source separation

此目录只用于保存**说明和清单**，不存放上游源码。

如需研究上游代码，运行：

```powershell
.\scripts\clone-references.ps1
```

默认会把第三方仓库克隆到本项目 Git 仓库之外：

```text
E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法\
├─ RimeTranslate\          # 你自己的 GitHub 仓库
├─ 第三方开源仓库\         # 上游源码，只供本地研究，不 push 到你的 repo
│  ├─ weasel\
│  ├─ rime-ice\
│  ├─ librime-lua\
│  ├─ ollama\
│  └─ rimes\
└─ 发布包\                 # 本地 Release 资产
```
