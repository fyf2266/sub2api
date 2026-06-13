# Sub2API 团队 Git 工作流与代码管理方案

> 🌿 由 Git Workflow Master 制定 | 适用于 Fork 上游仓库 + 频繁修改源码的团队协作场景

---

## 目录

1. [项目获取流程](#1-项目获取流程)
2. [部署方式对比与推荐](#2-部署方式对比与推荐)
3. [团队协作与代码管理流程](#3-团队协作与代码管理流程)
4. [日常操作速查手册](#4-日常操作速查手册)
5. [冲突处理与应急方案](#5-冲突处理与应急方案)

> 📖 **配套文档**：[日常实操手册](sub2api-daily-ops.md) | [混合部署详细实施手册](sub2api-deploy-guide.md) | [GitHub 认证排查指南](github-auth-troubleshooting.md)

---

## 1. 项目获取流程

### 1.1 整体策略：Fork + Upstream 同步

```
┌─────────────────────────────────────────────────────────────────┐
│                        GitHub 远程仓库                          │
│                                                                 │
│   Wei-Shaw/sub2api (upstream)     YourOrg/sub2api (origin)     │
│        │ 原始上游仓库                    │ 团队 Fork 仓库       │
│        │                                 │                      │
│        │    ┌──── fork ────►             │                      │
│        │    │                            │                      │
│        ▼    │                            ▼                      │
│   ┌─────────┴──────────┐       ┌────────────────┐              │
│   │  上游持续更新       │       │  团队自定义修改   │              │
│   │  (新功能/修复)      │       │  (定制功能)      │              │
│   └────────────────────┘       └────────────────┘              │
└─────────────────────────────────────────────────────────────────┘
                          │
                          ▼
                  ┌───────────────┐
                  │   本地开发机    │
                  │  (clone fork)  │
                  └───────────────┘
```

### 1.2 一步步操作

#### Step 1：在 GitHub 上 Fork 仓库

```
1. 打开 https://github.com/Wei-Shaw/sub2api
2. 点击右上角 "Fork" 按钮
3. 选择你的 GitHub 账号或组织作为 Fork 目标
4. 确认 Fork 完成（得到 YourOrg/sub2api 仓库）
```

#### Step 2：克隆 Fork 仓库到本地

```bash
# 克隆你的 Fork 仓库（不是原始仓库）
# 方式 A：HTTPS（推荐配合 PAT 使用）
git clone https://github.com/fyf2266/sub2api.git

# 方式 B：SSH（推荐长期使用，需先配置 SSH 密钥）
git clone git@github.com:fyf2266/sub2api.git

cd sub2api

# 验证当前远程仓库
git remote -v
# origin  https://github.com/fyf2266/sub2api.git (fetch)
# origin  https://github.com/fyf2266/sub2api.git (push)
```

> ⚠️ **重要提醒**：GitHub 已不支持密码认证！`git push` 时不能使用 GitHub 账户密码，必须使用 **Personal Access Token (PAT)** 或 **SSH 密钥**。详见 [GitHub 认证配置指南](github-auth-troubleshooting.md)

#### Step 3：配置认证（必须）

**快速方案：使用 PAT + 凭据缓存**

```bash
# 1. 在 GitHub 生成 PAT：Settings → Developer settings → Personal access tokens → Generate new token (classic)
#    勾选 repo 权限，复制生成的 Token

# 2. 配置凭据缓存（避免每次输入 Token）
git config --global credential.helper 'cache --timeout=86400'   # 缓存 24 小时
# 或永久保存（Token 明文存储，仅限个人服务器）
# git config --global credential.helper store

# 3. 首次推送时输入 Token 作为密码
git push origin main
# Username: fyf2266
# Password: ghp_xxxxxxxxxxxx   ← 粘贴 PAT（不是 GitHub 密码！）
```

**长期方案：使用 SSH 密钥（更安全，推荐服务器使用）**

```bash
# 1. 生成 SSH 密钥
ssh-keygen -t ed25519 -C "fyf2266@github" -f ~/.ssh/id_ed25519_github

# 2. 添加公钥到 GitHub：Settings → SSH and GPG keys → New SSH key
cat ~/.ssh/id_ed25519_github.pub   # 复制公钥内容

# 3. 切换 remote URL 为 SSH
git remote set-url origin git@github.com:fyf2266/sub2api.git

# 4. 测试连接
ssh -T git@github.com
```

> 📖 完整排查步骤和故障处理，请参考 [github-auth-troubleshooting.md](github-auth-troubleshooting.md)

#### Step 4：添加上游仓库

```bash
# 添加原始仓库作为 upstream
git remote add upstream https://github.com/Wei-Shaw/sub2api.git

# 禁止直接推送到上游（安全措施）
git remote set-url --push upstream no_push

# 验证远程仓库配置
git remote -v
# origin     https://github.com/fyf2266/sub2api.git (fetch)   ← 你的 Fork
# origin     https://github.com/fyf2266/sub2api.git (push)
# upstream   https://github.com/Wei-Shaw/sub2api.git (fetch)  ← 原始上游
# upstream   no_push (push)  ← 安全：防止误推到上游
```

#### Step 5：初始同步

```bash
# 拉取上游最新代码
git fetch upstream

# 确保本地 main 与上游 main 同步
git checkout main
git merge upstream/main

# 推送到你的 Fork 仓库
git push origin main
```

---

## 2. 部署方式对比与推荐

### 2.1 三种部署方式对比

| 维度 | 一键脚本安装 | Docker Compose | 源码编译部署 |
|------|-------------|----------------|-------------|
| **部署难度** | ⭐ 最简单 | ⭐⭐ 中等 | ⭐⭐⭐ 较复杂 |
| **环境要求** | 需自备 PostgreSQL + Redis | 仅需 Docker | 需 Go + Node.js + PostgreSQL + Redis |
| **升级方式** | 后台一键更新 | `docker compose pull && up -d` | `git pull` 重新编译 |
| **源码修改** | ❌ 不支持 | ⚠️ 需自建镜像 | ✅ 完全支持 |
| **迁移难度** | 中等 | ✅ 本地目录版极简 | 中等 |
| **资源占用** | 最低 | 中等（容器开销） | 最低 |
| **可观测性** | systemd 日志 | docker logs | 直接控制 |
| **调试能力** | 困难 | 中等 | ✅ 最强 |
| **热重载开发** | ❌ | ⚠️ 需挂载卷 | ✅ 前后端热重载 |
| **数据库管理** | 需自行管理 | 自动管理 | 需自行管理 |
| **适用场景** | 快速体验/不修改代码 | 标准生产部署 | 开发/频繁修改源码 |

### 2.2 详细优缺点分析

#### 方式一：一键脚本安装

**优点：**
- 一条命令完成安装，自动创建 systemd 服务
- 后台界面支持一键更新和回滚
- 资源占用最少，无容器开销
- 适合纯使用场景

**缺点：**
- 下载的是预编译二进制，**无法修改源码**
- 升级依赖官方发布版本，自定义修改会被覆盖
- 需要自行安装和管理 PostgreSQL、Redis
- 调试困难，无法热重载

#### 方式二：Docker Compose 部署

**优点：**
- 环境隔离，一键启动全部依赖（PostgreSQL、Redis）
- 本地目录版迁移极简（打包目录即可）
- 配置通过 `.env` 文件管理，清晰直观
- 社区活跃，文档完善

**缺点：**
- 默认拉取官方预构建镜像 `weishaw/sub2api:latest`，**无法直接修改源码**
- 如需自定义，必须修改源码后**自建 Docker 镜像**，增加构建步骤
- 容器有一定资源开销（约额外 100-200MB 内存）
- 调试不如源码编译方便

#### 方式三：源码编译部署

**优点：**
- ✅ **完全控制源码**，可随时修改、调试
- ✅ 前后端支持热重载开发（`go run` + `pnpm dev`）
- ✅ 修改后即时生效，无需重新构建镜像
- ✅ 可使用 Git 直接管理代码变更
- ✅ 调试能力最强，可直接使用 IDE 断点调试

**缺点：**
- 需要手动安装所有依赖（Go 1.25.7+, Node.js 18+, PostgreSQL 15+, Redis 7+）
- 编译步骤较多（前端 → 后端 → 嵌入）
- 需要自行管理 systemd 服务或进程守护
- 环境配置相对复杂

### 2.3 🏆 推荐：Docker Compose + 源码自建镜像（混合方案）

**为什么推荐这个方案？**

对于「个人服务器 + 频繁修改源码」的场景，纯源码编译虽然灵活但维护负担大，纯 Docker 部署虽然简单但无法修改源码。**混合方案**兼顾两者优势：

```
┌─────────────────────────────────────────────────────┐
│                    混合部署架构                       │
│                                                     │
│   源码仓库 (Git 管理)                                │
│       │                                             │
│       ▼                                             │
│   本地修改 → docker build 自建镜像                    │
│       │                                             │
│       ▼                                             │
│   Docker Compose 运行自建镜像                        │
│   (含 PostgreSQL + Redis + sub2api)                  │
│                                                     │
│   ✅ 源码可控    ✅ 环境隔离    ✅ 易迁移              │
│   ✅ 热重载开发  ✅ 生产部署一致                       │
└─────────────────────────────────────────────────────┘
```

**具体实施：详见 [sub2api-deploy-guide.md](sub2api-deploy-guide.md)**

> 📖 完整的实施手册已单独整理，包含：服务器环境准备、开发模式即时编译、生产模式自建镜像部署、一键脚本、常见问题排查等。

以下是核心流程速览：

#### 开发阶段：源码编译 + 热重载

```bash
# 1. 仅启动 PostgreSQL + Redis（Docker 容器）
cd deploy
docker compose -f docker-compose.local.yml -f docker-compose.dev-override.yml up -d postgres redis

# 2. 后端热重载（终端 1）— 使用 Air 自动重编译
cd backend && air
# 或直接: go run ./cmd/server/

# 3. 前端热重载（终端 2）— 修改即刷新
cd frontend && pnpm dev

# 浏览器访问 http://<IP>:5173 (前端) → API 代理到 :8080 (后端)
```

#### 部署阶段：自建 Docker 镜像

```bash
# 1. 构建自建镜像（多阶段构建：前端→后端→嵌入→运行时）
docker build -t sub2api-custom:latest \
    --build-arg GOPROXY=https://goproxy.cn,direct \
    -f Dockerfile .

# 2. 创建自定义 compose 文件（首次）
cd deploy
cp docker-compose.local.yml docker-compose.custom.yml
sed -i 's|image: weishaw/sub2api:latest|image: sub2api-custom:latest|' docker-compose.custom.yml

# 3. 启动服务
docker compose -f docker-compose.custom.yml up -d
```

#### 一键脚本（推荐）

```bash
# 一键构建部署（含备份、测试、构建、部署）
bash scripts/build-and-deploy.sh

# 同步上游 + 重新构建部署
bash scripts/sync-and-rebuild.sh

# 开发模式切换
bash scripts/dev-mode.sh start   # 切换到开发模式
bash scripts/dev-mode.sh stop    # 切换回生产模式
bash scripts/dev-mode.sh status  # 查看当前状态
```

---

## 3. 团队协作与代码管理流程

### 3.1 分支策略：Fork + Trunk-Based 混合模型

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          分支模型全景图                                   │
│                                                                          │
│  upstream/main  ●────●────●────●────●────●────●────●────●────●         │
│                  │         │         │         │         │               │
│  origin/main    ●────●────●────●────●────●────●────●────●────●         │
│                  │    │    │    │    │    │    │    │    │               │
│                  │  sync   │  sync   │  sync   │  sync   │              │
│                  │         │         │         │         │               │
│  feat/xxx       └──●──●──┘                        │         │           │
│  fix/yyy                  └──●──┘                  │         │           │
│  feat/zzz                          └──●──●──●──┘  │         │           │
│  custom/xxx                                 └──●──┘│         │           │
│                                                    │         │           │
│  ◆ main = 始终与 upstream 保持同步 + 团队自定义基础                    │
│  ◆ feat/* = 功能开发分支（短生命周期）                                   │
│  ◆ fix/* = 修复分支（短生命周期）                                       │
│  ◆ custom/* = 自定义功能分支（可稍长但尽量短）                           │
└──────────────────────────────────────────────────────────────────────────┘
```

### 3.2 核心原则

| 原则 | 说明 |
|------|------|
| **main 始终可同步** | `origin/main` 是上游同步的基础，永远保持可以 `merge upstream/main` 的状态 |
| **自定义走分支** | 所有修改都在功能分支完成，不直接在 main 上改 |
| **冲突在分支解决** | 同步上游时产生冲突，在功能分支解决，不污染 main |
| **Conventional Commits** | 使用 `feat:`/`fix:`/`custom:`/`chore:` 前缀 |
| **原子提交** | 每个提交只做一件事，可独立回滚 |

### 3.3 完整工作流

#### 流程一：同步上游更新（核心流程，定期执行）

```bash
# ============================================================
# 🔄 同步上游更新 — 建议每天或每周定期执行
# ============================================================

# Step 1: 获取上游最新代码
git fetch upstream

# Step 2: 切换到 main 并同步
git checkout main
git merge upstream/main

# Step 3: 解决冲突（如有）
# 冲突文件会标记出来，手动解决后：
git add .
git commit -m "chore: sync upstream Wei-Shaw/sub2api"

# Step 4: 推送到团队 Fork
git push origin main

# Step 5: 同步到所有进行中的功能分支
git checkout feat/my-feature
git rebase origin/main   # 推荐 rebase，保持线性历史
# 如果冲突，解决后继续：
git rebase --continue
git push --force-with-lease origin feat/my-feature
```

#### 流程二：开发自定义功能

```bash
# ============================================================
# ✨ 开发自定义功能 — 从最新 main 创建分支
# ============================================================

# Step 1: 确保从最新的 main 创建分支
git fetch origin
git checkout -b custom/my-feature origin/main

# Step 2: 开发过程中，遵循原子提交
# 修改代码...
git add -p              # 精确选择要提交的更改
git commit -m "custom: add xxx functionality"

# Step 3: 开发期间定期同步上游（避免积压大量冲突）
git fetch upstream
git rebase upstream/main

# Step 4: 开发完成，推送到 Fork 仓库
git push origin custom/my-feature

# Step 5: 在 GitHub 上创建 PR → origin/main
# （团队 Code Review 后合并）
```

#### 流程三：合并自定义功能到 main

```bash
# ============================================================
# 🔀 合并功能到 main — 通过 PR 或本地合并
# ============================================================

# 方式 A：通过 GitHub PR 合并（推荐，有 Code Review 记录）

# 方式 B：本地合并（适合小团队快速迭代）
git checkout main
git merge --no-ff custom/my-feature
# --no-ff 保留分支合并记录，方便回溯
git push origin main

# 清理已合并的分支
git branch -d custom/my-feature
git push origin --delete custom/my-feature
```

#### 流程四：紧急修复

```bash
# ============================================================
# 🚑 紧急修复 — 基于当前 main 快速修复
# ============================================================

git checkout -b fix/hotfix-xxx origin/main
# 修复代码...
git commit -m "fix: resolve xxx issue"
git push origin fix/hotfix-xxx

# 合并到 main
git checkout main
git merge --no-ff fix/hotfix-xxx
git push origin main

# 同步到上游（如果修复是通用的，可以提交 PR 到上游）
git push upstream fix/hotfix-xxx
# 然后在 GitHub 上向上游创建 PR
```

### 3.4 自定义代码管理的关键规则

#### 规则一：区分「自定义修改」与「上游兼容修改」

```
┌─────────────────────────────────────────────────────┐
│           修改类型与处理策略                           │
│                                                     │
│  ┌─────────────────┐    ┌──────────────────────┐    │
│  │ 自定义修改       │    │ 上游兼容修改          │    │
│  │ (custom: 前缀)  │    │ (feat:/fix: 前缀)    │    │
│  ├─────────────────┤    ├──────────────────────┤    │
│  │ • 仅团队需要    │    │ • 对上游也有价值      │    │
│  │ • 不回推上游    │    │ • 可提交 PR 到上游    │    │
│  │ • 独立维护      │    │ • 减少长期冲突        │    │
│  └─────────────────┘    └──────────────────────┘    │
│                                                     │
│  💡 尽量将修改设计为「上游兼容」，减少长期维护负担      │
└─────────────────────────────────────────────────────┘
```

#### 规则二：避免冲突的代码组织方式

**最小侵入原则 — 优先使用配置而非修改源码：**

| 方式 | 侵入性 | 示例 |
|------|--------|------|
| 环境变量 / 配置文件 | ⭐ 最低 | 修改 `.env` 或 `config.yaml` |
| 新增独立文件 | ⭐⭐ 低 | 新增 handler/service，不修改已有文件 |
| 扩展已有模块 | ⭐⭐⭐ 中 | 在已有文件中新增方法/函数 |
| 修改核心逻辑 | ⭐⭐⭐⭐ 高 | 修改上游已有代码，冲突风险大 |
| 重构上游结构 | ⭐⭐⭐⭐⭐ 最高 | 尽量避免 |

**实践建议：**
1. **新功能放在独立文件**：新增 `custom_*.go` 文件，而非修改上游已有文件
2. **修改点集中管理**：如果必须修改上游文件，尽量将修改集中到少量位置
3. **使用 Git 记录修改清单**：维护一个 `CUSTOM_CHANGES.md` 记录所有自定义修改位置

#### 规则三：维护自定义修改清单

创建 `CUSTOM_CHANGES.md` 文件追踪所有自定义修改：

```markdown
# 自定义修改清单

## 修改文件列表

| 文件路径 | 修改类型 | 修改说明 | 同步风险 |
|----------|----------|----------|----------|
| backend/internal/handler/custom_xxx.go | 新增 | 自定义XXX功能 | 无 |
| backend/internal/service/user.go:45 | 修改 | 添加自定义验证逻辑 | 中 |
| frontend/src/views/Dashboard.vue:120 | 修改 | 自定义仪表盘布局 | 高 |
```

每次同步上游时，优先检查这些文件的变更。

### 3.5 CI/CD 集成建议

```
┌────────────────────────────────────────────────────────────┐
│                     CI/CD 流水线                           │
│                                                            │
│  Push to branch ──► GitHub Actions ──► 检查结果            │
│                         │                                  │
│                    ┌────┴────┐                              │
│                    ▼         ▼                              │
│               后端检查    前端检查                           │
│               │          │                                  │
│               ▼          ▼                                  │
│          go test    pnpm build                             │
│          golangci   pnpm lint                              │
│               │          │                                  │
│               └────┬─────┘                                 │
│                    ▼                                       │
│              合并到 main?                                   │
│              │         │                                    │
│             Yes        No                                  │
│              │         │                                    │
│              ▼         返回修改                             │
│      自动构建镜像                                         │
│      docker build                                        │
│              │                                             │
│              ▼                                            │
│      部署到服务器                                         │
│      (手动/自动)                                          │
└────────────────────────────────────────────────────────────┘
```

**GitHub Actions 工作流配置（`.github/workflows/custom-ci.yml`）：**

```yaml
name: Custom CI

on:
  push:
    branches: [main, 'custom/**', 'feat/**', 'fix/**']
  pull_request:
    branches: [main]

jobs:
  backend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25.7'
      - name: Run unit tests
        working-directory: backend
        run: go test -tags=unit ./...
      - name: Run linter
        working-directory: backend
        run: |
          go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7
          golangci-lint run ./...

  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v4
        with:
          version: 9
      - uses: actions/setup-node@v4
        with:
          node-version: 18
          cache: pnpm
          cache-dependency-path: frontend/pnpm-lock.yaml
      - name: Install dependencies
        working-directory: frontend
        run: pnpm install --frozen-lockfile
      - name: Build
        working-directory: frontend
        run: pnpm run build
```

---

## 4. 日常操作速查手册

### 4.1 同步上游（每日/每周）

```bash
# 快速同步
git fetch upstream
git checkout main && git merge upstream/main && git push origin main
```

### 4.2 开发新功能

```bash
# 创建功能分支
git checkout -b custom/feature-name origin/main

# 开发 + 提交
git add -p && git commit -m "custom: description"

# 推送
git push origin custom/feature-name
```

### 4.3 同步到功能分支

```bash
# 定期同步（在功能分支开发期间）
git fetch upstream
git rebase upstream/main
git push --force-with-lease origin custom/feature-name
```

### 4.4 构建部署

```bash
# 开发模式
cd backend && go run ./cmd/server/     # 后端
cd frontend && pnpm dev                # 前端

# 生产构建（自建 Docker 镜像）
docker build -t sub2api-custom:latest .
docker compose -f docker-compose.local.yml up -d
```

### 4.5 修改 Ent Schema 后

```bash
cd backend
go generate ./ent          # 重新生成 Ent 代码
go generate ./cmd/server   # 重新生成 Wire 代码
git add ent/               # 生成的文件也要提交
```

---

## 5. 冲突处理与应急方案

### 5.1 同步上游时的冲突处理

```bash
# ============================================================
# ⚔️ 同步上游时遇到冲突
# ============================================================

git checkout main
git merge upstream/main
# 如果有冲突：

# Step 1: 查看冲突文件
git status

# Step 2: 逐个文件解决冲突
# 打开冲突文件，选择保留的内容：
# <<<<<<< HEAD         ← 你的修改
# =======
# >>>>>>> upstream/main ← 上游的修改

# Step 3: 标记为已解决
git add <resolved-file>

# Step 4: 完成合并
git commit -m "chore: sync upstream, resolve conflicts"

# Step 5: 推送
git push origin main
```

### 5.2 Rebase 冲突处理

```bash
# ============================================================
# ⚔️ Rebase 功能分支时遇到冲突
# ============================================================

git rebase upstream/main
# 如果有冲突：

# Step 1: 解决冲突文件
# Step 2: 标记为已解决
git add <resolved-file>

# Step 3: 继续 rebase
git rebase --continue

# Step 4: 安全推送（覆盖远程分支）
git push --force-with-lease origin custom/my-feature

# 🔴 如果冲突太复杂，可以放弃 rebase：
git rebase --abort
# 改用 merge 策略：
git merge upstream/main
```

### 5.3 回滚方案

```bash
# ============================================================
# 🔄 回滚 — 各种场景的恢复方法
# ============================================================

# 场景 1：最近的提交有问题，回退一个提交
git revert HEAD
git push origin main

# 场景 2：合并上游后有严重问题
git log --oneline -5          # 找到合并前的提交
git revert -m 1 <merge-commit>   # 回退合并
git push origin main

# 场景 3：功能分支搞砸了，重新开始
git checkout main
git branch -D custom/broken-feature   # 删除本地分支
git push origin --delete custom/broken-feature  # 删除远程分支
git checkout -b custom/new-start origin/main    # 重新创建

# 场景 4：docker 部署出问题，回退到上一版本
docker compose -f docker-compose.local.yml down
docker tag sub2api-custom:latest sub2api-custom:backup
docker build -t sub2api-custom:previous <previous-commit> .
docker compose -f docker-compose.local.yml up -d
```

### 5.4 reflog 急救

```bash
# ============================================================
# 🆘 reflog — Git 的"后悔药"
# ============================================================

# 如果误操作丢失了提交，用 reflog 找回：
git reflog

# 输出类似：
# abc1234 HEAD@{0}: commit: custom: add feature
# def5678 HEAD@{1}: checkout: moving from main to custom/xxx
# 901abcd HEAD@{2}: merge upstream/main: Fast-forward

# 找到想要的提交，恢复：
git reset --hard abc1234
# 或创建新分支保存：
git branch recovery-branch abc1234
```

### 5.5 定期维护

```bash
# ============================================================
# 🧹 定期仓库维护
# ============================================================

# 清理已合并的本地分支
git branch --merged main | grep -v "main" | xargs -r git branch -d

# 清理远程已删除的分支引用
git remote prune origin

# 清理 Docker 旧镜像
docker image prune -f

# 查看仓库健康状态
git fsck
```

---

## 附录 A：Conventional Commits 规范

| 前缀 | 用途 | 示例 |
|------|------|------|
| `feat:` | 新功能 | `feat: add user quota management` |
| `fix:` | 修复 Bug | `fix: resolve session timeout issue` |
| `custom:` | 团队自定义功能 | `custom: add billing dashboard` |
| `chore:` | 杂务（构建/依赖） | `chore: update dependencies` |
| `docs:` | 文档 | `docs: update deployment guide` |
| `refactor:` | 重构 | `refactor: simplify auth logic` |
| `test:` | 测试 | `test: add unit tests for gateway` |
| `style:` | 代码风格 | `style: format with gofmt` |
| `sync:` | 同步上游 | `sync: merge upstream v1.2.0` |

## 附录 B：Git Remote 配置参考

```bash
# 查看当前远程仓库
git remote -v

# 标准配置：
# origin   → 你的 Fork 仓库（推送/拉取）
# upstream → 原始上游仓库（仅拉取）

# 禁止直接推送到上游（防止误操作）
git remote set-url --push upstream no_push

# 验证配置
git remote -v
# origin     https://github.com/YourOrg/sub2api.git (fetch)
# origin     https://github.com/YourOrg/sub2api.git (push)
# upstream   https://github.com/Wei-Shaw/sub2api.git (fetch)
# upstream   no_push (push)  ← 安全措施
```

## 附录 C：一键同步脚本

创建 `scripts/sync-upstream.sh`：

```bash
#!/bin/bash
set -e

echo "🔄 同步上游仓库 Wei-Shaw/sub2api ..."

# 获取上游最新
git fetch upstream

# 切换到 main
git checkout main

# 合并上游
echo "📦 合并 upstream/main ..."
if git merge upstream/main; then
    echo "✅ 合并成功，无冲突"
else
    echo "⚠️ 存在冲突，请手动解决后执行："
    echo "   git add ."
    echo "   git commit -m 'sync: merge upstream'"
    echo "   git push origin main"
    exit 1
fi

# 推送到 Fork
git push origin main

echo "🎉 同步完成！"

# 列出进行中的功能分支，提醒同步
BRANCHES=$(git branch --list 'custom/*' 'feat/*' 'fix/*' --format='%(refname:short)')
if [ -n "$BRANCHES" ]; then
    echo ""
    echo "📋 以下功能分支可能需要 rebase："
    echo "$BRANCHES" | while read branch; do
        echo "   git checkout $branch && git rebase origin/main"
    done
fi
```
