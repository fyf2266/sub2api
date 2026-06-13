# Git 仓库层级结构与上游同步指南

## 目录
- [Git 仓库层级结构](#1-git-仓库层级结构)
- [上游同步流程](#2-上游同步流程)
- [同步命令详解](#3-同步命令详解)
- [同步流向](#4-同步流向)
- [推荐工作流程](#5-推荐工作流程)
- [当前状态检查](#6-当前状态检查)

---

## 1. Git 仓库层级结构

### 1.1 三个层级概览

```
┌─────────────────────────────────────────────────────────────────┐
│                     Git 仓库层级结构                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1️⃣ 远程仓库（Remote） - GitHub 服务器                            │
│     ├── origin (你的仓库)                                        │
│     │   ├── URL: https://github.com/fyf2266/sub2api.git         │
│     │   └── 权限: 可以 push/pull                                 │
│     │                                                            │
│     └── upstream (上游仓库)                                      │
│         ├── URL: https://github.com/Wei-Shaw/sub2api.git        │
│         └── 权限: 只能 fetch（读取），不能 push                    │
│                                                                  │
│  2️⃣ 远程引用（Remote References） - 本地缓存的远程分支             │
│     ├── origin/main      → 你的远程 main 分支                    │
│     ├── origin/cla-signatures                                   │
│     ├── upstream/main    → 上游的 main 分支                      │
│     ├── upstream/cla-signatures                                 │
│     ├── upstream/dev                                           │
│     └── upstream/preview                                       │
│     ⚠️ 这些只是"指针"，不会影响你的本地代码                        │
│                                                                  │
│  3️⃣ 本地分支（Local Branches） - 你实际工作的分支                  │
│     ├── main              → 追踪 origin/main                     │
│     └── cla-signatures    → 追踪 upstream/cla-signatures         │
│     ✅ 这些是你真正修改和提交代码的地方                             │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 各层级详细说明

#### 第一层：远程仓库（Remote）

| 远程名称 | URL | 说明 |
|---------|-----|------|
| `origin` | `https://github.com/fyf2266/sub2api.git` | 你的 GitHub 仓库，可读写 |
| `upstream` | `https://github.com/Wei-Shaw/sub2api.git` | 上游仓库，**只读**（设置为 no_push） |

#### 第二层：远程引用（Remote References）

这些是本地存储的远程分支"指针"，类似于书签：

| 引用名称 | 指向 | 用途 |
|---------|------|------|
| `origin/main` | 你的远程 main | 追踪你的仓库主分支 |
| `origin/cla-signatures` | 你的远程 CLA 分支 | （可能不存在） |
| `upstream/main` | 上游的 main | 获取上游最新代码 |
| `upstream/cla-signatures` | 上游的 CLA 分支 | 获取 CLA 签名记录 |
| `upstream/dev` | 上游的开发分支 | 可能的开发版本 |
| `upstream/preview` | 上游的预览分支 | 预览版本 |

#### 第三层：本地分支（Local Branches）

你实际工作和提交代码的地方：

| 分支名 | 追踪 | 用途 |
|--------|------|------|
| `main` | `origin/main` | 你的主开发分支 |
| `cla-signatures` | `upstream/cla-signatures` | CLA 签名同步分支 |

---

## 2. 上游同步流程

### 2.1 同步流程图

```
上游仓库 (Wei-Shaw/sub2api)
      │
      │ git fetch upstream
      ▼
┌────────────────────────────────────────┐
│  本地远程引用更新                        │
│  upstream/main, upstream/cla-signatures │
│  upstream/dev, upstream/preview 等      │
│  ✅ 只是更新"指针"，不改变你的代码         │
└────────────────────────────────────────┘
      │
      │ git merge / git rebase
      ▼
┌────────────────────────────────────────┐
│  本地分支更新                           │
│  main, cla-signatures 等               │
│  ✅ 这才会真正更新你的代码               │
└────────────────────────────────────────┘
```

### 2.2 两步同步法

#### 第一步：获取上游更新（更新远程引用）

```bash
# 获取所有上游远程分支的更新
git fetch upstream

# 获取特定分支
git fetch upstream main
git fetch upstream cla-signatures
```

#### 第二步：同步到本地分支

**同步 cla-signatures 分支：**
```bash
git checkout cla-signatures
git merge upstream/cla-signatures
```

**同步 main 分支：**
```bash
git checkout main
git merge upstream/main
```

---

## 3. 同步命令详解

### 3.1 git fetch

**作用**：从远程仓库下载最新提交，但**不修改**你的本地分支

```bash
git fetch upstream        # 获取上游所有分支更新
git fetch upstream main   # 只获取上游 main 分支
git fetch --all          # 获取所有远程仓库更新
```

**执行后**：更新 `upstream/*` 远程引用，本地代码不变

### 3.2 git merge

**作用**：将上游分支合并到当前本地分支

```bash
git merge upstream/main          # 合并（可能产生合并提交）
git merge --no-ff upstream/main  # 禁止快进，保留分支历史
```

**特点**：
- 保留完整的历史记录
- 可能产生合并提交
- 适合协作开发场景

### 3.3 git rebase（推荐）

**作用**：将你的提交"变基"到上游分支之上

```bash
git rebase upstream/main          # 变基（保持历史线性）
git rebase -i upstream/main      # 交互式变基，可以修改提交
```

**特点**：
- 历史线性整洁
- 没有合并提交
- 适合个人开发或同步上游代码

**注意**：⚠️ **不要对已推送到远程的提交进行变基！**

---

## 4. 同步流向

### 4.1 数据流向图

```
同步内容流向：

upstream/main (远程引用)
    ↓ merge/rebase
main (本地分支)
    ↓ push
origin/main (你的远程仓库)

upstream/cla-signatures (远程引用)
    ↓ merge/rebase
cla-signatures (本地分支)
    ⚠️ 不能推送到 origin（因为这是 CLA 签名分支）
```

### 4.2 各分支同步策略

| 本地分支 | 上游来源 | 同步方式 | 推送到 origin？ |
|---------|---------|---------|---------------|
| `main` | `upstream/main` | fetch + merge/rebase | ✅ 是 |
| `cla-signatures` | `upstream/cla-signatures` | fetch + merge | ❌ 否（不需要） |

### 4.3 实际同步示例

#### 场景 1：同步 main 分支

```bash
# 1. 确保在 main 分支
git checkout main

# 2. 获取上游更新
git fetch upstream

# 3. 合并上游 main 到本地 main
git merge upstream/main

# 4. 推送到你的远程仓库
git push origin main
```

#### 场景 2：同步 cla-signatures 分支

```bash
# 1. 切换到 cla-signatures
git checkout cla-signatures

# 2. 获取上游更新
git fetch upstream

# 3. 合并上游的 CLA 签名
git merge upstream/cla-signatures

# 4. ⚠️ 注意：不能也不需要推送到 origin
```

---

## 5. 推荐工作流程

### 5.1 日常同步上游代码

```bash
# 1. 获取上游所有更新
git fetch upstream

# 2. 切换到 main 分支
git checkout main

# 3. 合并上游 main
git merge upstream/main

# 4. 推送到你的仓库
git push origin main
```

### 5.2 偶尔同步 CLA 签名（保持记录）

```bash
# 1. 切换到 cla-signatures
git checkout cla-signatures

# 2. 获取上游更新
git fetch upstream

# 3. 合并上游的 CLA 签名
git merge upstream/cla-signatures

# 4. ⚠️ 不要 push，这个分支不需要同步到 origin
```

### 5.3 完整开发流程

```bash
# 1. 同步上游代码
git fetch upstream
git checkout main
git merge upstream/main

# 2. 创建功能分支
git checkout -b feature/your-feature

# 3. 开发并提交
git add .
git commit -m "feat: 添加新功能"

# 4. 推送功能分支到你的仓库
git push -u origin feature/your-feature

# 5. 开发完成后合并到 main
git checkout main
git merge feature/your-feature
git push origin main
```

---

## 6. 当前状态检查

### 6.1 查看各分支与上游的差异

```bash
# main 分支落后多少（上游有新提交）
git log main..upstream/main

# main 分支领先多少（本地有新提交）
git log upstream/main..main

# cla-signatures 落后多少
git log cla-signatures..upstream/cla-signatures
```

### 6.2 查看当前分支状态

```bash
# 查看本地分支及其追踪关系
git branch -vv

# 查看所有远程分支
git branch -r

# 查看当前工作区状态
git status
```

### 6.3 当前环境信息

- **当前分支**：`main`
- **追踪关系**：`main` → `origin/main`
- **upstream 追踪**：`upstream/main`, `upstream/cla-signatures`
- **upstream 状态**：已同步，无差异

---

## 常见问题

### Q1: 为什么 cla-signatures 和 main 没有共同祖先？
**A**: 这是两个**完全独立的分支**，来自不同的仓库工作流：
- `main` 是正常的代码开发分支
- `cla-signatures` 是上游仓库专门用于收集 CLA 签名的特殊分支
- 它们不需要合并，只是在本地各自追踪不同的远程分支

### Q2: 可以把 cla-signatures 推送到 origin 吗？
**A**: 不建议。`cla-signatures` 是上游仓库的特殊分支，你的 origin 仓库不需要这个分支。

### Q3: 同步时出现冲突怎么办？
**A**: 
1. 先备份：`git branch backup-branch`
2. 查看冲突文件：`git status`
3. 手动解决冲突
4. 添加解决后的文件：`git add <file>`
5. 完成合并：`git commit` 或 `git rebase --continue`

---

## 相关命令速查表

| 操作 | 命令 |
|------|------|
| 查看远程仓库 | `git remote -v` |
| 获取上游更新 | `git fetch upstream` |
| 合并上游分支 | `git merge upstream/main` |
| 变基到上游 | `git rebase upstream/main` |
| 推送到你的仓库 | `git push origin main` |
| 查看分支差异 | `git log main..upstream/main` |
| 查看追踪关系 | `git branch -vv` |
| 创建功能分支 | `git checkout -b feature/xxx` |
| 切换分支 | `git checkout <branch>` |

---

*文档最后更新：2026-06-11*
