# GitHub 发布操作指南（按当前本地路径）

你的工作区：

```text
E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法
```

推荐本地结构：

```text
实用类——双语输入法\
├─ RimeTranslate\          # 你自己的 GitHub 仓库
├─ 第三方开源仓库\         # 上游源码，仅供研究，不 push
└─ 发布包\                 # 本地 Release 资产
```

## 1. 创建目录

```powershell
$root = 'E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法'
New-Item -ItemType Directory -Force `
  (Join-Path $root 'RimeTranslate'), `
  (Join-Path $root '第三方开源仓库'), `
  (Join-Path $root '发布包') | Out-Null
```

将本源码包解压后的全部文件放入：

```text
E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法\RimeTranslate
```

## 2. 初始化 Git

```powershell
cd 'E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法\RimeTranslate'

git init
git branch -M main
git add .
git commit -m "Initial release: RimeTranslate v2.2.0"
```

## 3. GitHub 上创建空仓库

建议仓库名：

```text
RimeTranslate
```

创建时不要勾选自动生成 README / License / .gitignore，因为本地已经有。

然后把 GitHub 给你的仓库地址替换到下面：

```powershell
git remote add origin https://github.com/YOUR_USERNAME/RimeTranslate.git
git push -u origin main
```

## 4. 发布 v2.2.0

```powershell
git tag -a v2.2.0 -m "RimeTranslate v2.2.0"
git push origin v2.2.0
```

GitHub Actions 会自动生成 Windows x64 Release 包。

## 5. 第三方仓库单独存放

如果只是想本地研究上游代码：

```powershell
.\scripts\clone-references.ps1
```

它们会被克隆到：

```text
...\实用类——双语输入法\第三方开源仓库\
```

不会进入你自己的 Git 仓库。

## 6. 日后开发流程

```powershell
cd 'E:\桌面重要内容\AI\个人github项目管理\项目主目录\实用类——双语输入法\RimeTranslate'

git status
git add .
git commit -m "Describe your change"
git push
```

发布新版本时：

```powershell
git tag -a v2.3.0 -m "RimeTranslate v2.3.0"
git push origin v2.3.0
```
