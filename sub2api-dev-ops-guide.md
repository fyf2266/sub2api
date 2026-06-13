# Sub2API 开发环境操作手册

> 本文档记录 Docker Compose + 源码自建镜像（混合方案）的开发环境启动/关闭步骤及配置修改记录

---

## 一、服务启动步骤

### 0. 启动前检查与清理（重要！）

```bash
# 检查并关闭已运行的服务
cd /home/root/sub2api

echo "=== 检查并关闭已运行的服务 ==="

# 关闭前端进程（精确匹配）
echo "关闭前端进程..."
pkill -f "pnpm dev" 2>/dev/null || true
pkill -f "vite" 2>/dev/null || true

# 关闭后端进程（精确匹配，避免误杀系统进程）
echo "关闭后端进程..."
pkill -f "/root/go/bin/air" 2>/dev/null || true
pkill -f "go run ./cmd/server" 2>/dev/null || true
pkill -f "./backend/server" 2>/dev/null || true
pkill -f "./tmp/main" 2>/dev/null || true

# 等待进程终止
sleep 2

# 验证服务已关闭
echo ""
echo "=== 验证服务已关闭 ==="
ps aux | grep -E "pnpm dev|vite|/root/go/bin/air|go run ./cmd/server|./backend/server|./tmp/main" | grep -v grep || echo "✓ 所有服务已关闭"

# 检查端口占用
echo ""
echo "=== 检查端口占用 ==="
ss -tlnp | grep -E "8080|3000|3001|5173" || echo "✓ 所有端口已释放"
```

**⚠️ 注意：** 此步骤确保启动前没有残留的服务进程，避免端口冲突和重复启动问题。

### 1. 启动数据库容器（PostgreSQL + Redis）

```bash
cd /home/root/sub2api/deploy

# 启动 PostgreSQL 和 Redis 容器（开发模式，暴露端口）
docker compose -f docker-compose.local.yml -f docker-compose.dev-override.yml up -d postgres redis

# 验证容器状态
docker compose -f docker-compose.local.yml ps
```

**预期结果：**
- PostgreSQL: 端口 `5432`，状态 `healthy`
- Redis: 端口 `6379`，状态 `healthy`

### 2. 启动后端服务

**⚠️ 重要：启动后端前必须先加载环境变量！**

```bash
# 方式 A：开发调试（推荐，修改源码后立即生效）
cd /home/root/sub2api/deploy && set -a && source .env && set +a && cd ../backend && go run ./cmd/server/

# 方式 B：使用 air 热重载（自动重新编译）
# 安装 air: go install github.com/air-verse/air@latest
cd /home/root/sub2api/deploy && set -a && source .env && set +a && cd ../backend && $(go env GOPATH)/bin/air

# 方式 C：生产部署（需先编译，修改源码后需重新编译）
cd /home/root/sub2api/backend && go build -o server ./cmd/server/
cd /home/root/sub2api/deploy && set -a && source .env && set +a && cd ../backend && ./server
```

**关键区别：**
| 方式 | 适用场景 | 源码修改后 |
|------|----------|------------|
| `go run` | 开发调试 | 立即生效（重新运行） |
| `air` | 开发调试 | 自动重新编译运行 |
| `./server` | 生产部署 | 需手动重新编译 |

**验证后端启动：**
```bash
curl http://localhost:8080/health
# 预期返回: {"status":"ok"}
```

### 3. 启动前端服务

```bash
cd /home/root/sub2api/frontend

# 安装依赖（首次或依赖更新后）
pnpm install

# 启动开发服务器
pnpm dev
```

**访问前端：**
- 开发地址: `http://localhost:5173`
- API 自动代理到后端 `http://localhost:8080`

### 4. 一键启动脚本（推荐）

```bash
# 创建启动脚本
cd /home/root/sub2api

# 启动所有服务
bash -c '
# ===== 启动前清理 =====
echo "=== 关闭已有服务 ==="

# 关闭前端进程（精确匹配）
echo "关闭前端进程..."
pkill -f "pnpm dev" 2>/dev/null || true
pkill -f "vite" 2>/dev/null || true

# 关闭后端进程（精确匹配，避免误杀系统进程）
echo "关闭后端进程..."
pkill -f "/root/go/bin/air" 2>/dev/null || true
pkill -f "./backend/server" 2>/dev/null || true
pkill -f "./tmp/main" 2>/dev/null || true

sleep 2

# ===== 启动服务 =====
echo ""
echo "=== 启动数据库容器 ==="
cd deploy && docker compose -f docker-compose.local.yml -f docker-compose.dev-override.yml up -d postgres redis

# 等待数据库就绪
echo ""
echo "=== 等待数据库启动 (10秒) ==="
sleep 10

# 启动后端（air 热重载模式）
echo ""
echo "=== 启动后端服务 ==="
cd ../backend && $(go env GOPATH)/bin/air &

# 等待后端启动
sleep 5

# 启动前端
echo ""
echo "=== 启动前端服务 ==="
cd ../frontend && pnpm dev &

echo ""
echo "=== 服务启动完成 ==="
echo "前端地址: http://localhost:3000"
echo "后端地址: http://localhost:8080"
'
```

**⚠️ 注意：**
- 此脚本会先关闭所有已有服务，避免端口冲突
- 使用 `air` 启动后端，修改 Go 源码后自动重新编译运行
- 需要先创建 `.air.toml` 配置文件（见 2. 节）

---

## 二、服务关闭步骤

### 1. 关闭后端进程

```bash
pkill -f "backend/server"
# 或
pkill -f "server"
```

### 2. 关闭前端进程

```bash
pkill -f "pnpm dev"
```

### 3. 关闭 Docker 容器

```bash
cd /home/root/sub2api/deploy
docker compose -f docker-compose.local.yml down
```

### 4. 一键关闭脚本

```bash
# 关闭所有服务
bash -c '
pkill -f "backend/server"
pkill -f "pnpm dev"
cd /home/root/sub2api/deploy && docker compose -f docker-compose.local.yml down
'
```

### 5. 验证服务已关闭

```bash
# 检查 Docker 容器
docker compose -f docker-compose.local.yml ps
# 应返回空列表

# 检查进程
ps aux | grep -E "server|pnpm" | grep -v grep

# 检查端口
ss -tlnp | grep -E "8080|5173|5432|6379"
```

---

## 三、登录信息

| 项目 | 值 |
|------|-----|
| **管理员邮箱** | `admin@sub2api.local` |
| **管理员密码** | `Admin@123456` |
| **后端地址** | `http://localhost:8080` |
| **前端地址** | `http://localhost:5173` |

---

## 四、修改的文件记录

### 4.1 deploy/.env

**路径:** `/home/root/sub2api/deploy/.env`

**修改内容:**

| 行号 | 原值 | 新值 | 说明 |
|------|------|------|------|
| 109 | 无 | `DATABASE_PASSWORD=263fad824a6bae750f6ad5690cdd894f` | 新增：后端连接数据库的密码变量 |
| 183 | `ADMIN_PASSWORD=` | `ADMIN_PASSWORD=Admin@123456` | 设置管理员密码 |

**修改原因:**
- 后端程序读取 `DATABASE_PASSWORD` 环境变量连接数据库，而 Docker 容器使用 `POSTGRES_PASSWORD`
- 设置固定管理员密码，避免每次启动随机生成

### 4.2 backend/config.yaml

**路径:** `/home/root/sub2api/backend/config.yaml`

**修改内容:**

| 行号 | 原值 | 新值 | 说明 |
|------|------|------|------|
| 10 | `password: "your_postgres_password"` | `password: "263fad824a6bae750f6ad5690cdd894f"` | 更新数据库密码 |
| 14-15 | `max_open_conns: 10` / `max_idle_conns: 5` | `max_open_conns: 256` / `max_idle_conns: 128` | 更新连接池配置 |
| 26 | `secret: "your_jwt_secret"` | `secret: "d196446f4e2d48f504476eafd719686ee85803f912adcad12767879d2a2e480c"` | 更新 JWT 密钥 |

**修改原因:**
- config.yaml 中的密码需要与 .env 中的 `POSTGRES_PASSWORD` 一致
- JWT 密钥需要与 .env 中的 `JWT_SECRET` 一致
- 连接池配置需要匹配生产环境需求

### 4.3 数据库用户表

**操作:** 在 PostgreSQL 数据库中创建管理员用户

```sql
INSERT INTO users (email, password_hash, role, balance, concurrency, status, created_at, updated_at, username, notes, wechat, totp_enabled, balance_notify_enabled, total_recharged, signup_source, rpm_limit)
VALUES (
    'admin@sub2api.local',
    '$2a$10$M4zGw317oNkSaSNdf2Ra3uUp8PMp9iA3051qEG4/tw9H.Er8Ns/Ry',  -- bcrypt hash of "Admin@123456"
    'admin',
    0,
    30,
    'active',
    NOW(),
    NOW(),
    'admin',
    '',
    '',
    false,
    true,
    0,
    'email',
    0
);
```

---

## 五、服务架构图

```
┌─────────────────────────────────────────────────────┐
│                    开发模式架构                       │
│                                                     │
│   Docker 容器:                                       │
│   ├── PostgreSQL (postgres:18-alpine)              │
│   │   端口: 5432                                    │
│   │   数据目录: ./postgres_data                     │
│   │                                                 │
│   └── Redis (redis:8-alpine)                        │
│       端口: 6379                                    │
│       数据目录: ./redis_data                        │
│                                                     │
│   本地进程:                                          │
│   ├── 后端: ./server (端口 8080)                    │
│   │   配置: backend/config.yaml                     │
│   │                                                 │
│   └── 前端: pnpm dev (端口 5173)                    │
│       API 代理: → localhost:8080                    │
│                                                     │
│   环境变量: deploy/.env                              │
└─────────────────────────────────────────────────────┘
```

---

## 六、常用命令速查

### 服务状态检查

```bash
# Docker 容器状态
docker compose -f docker-compose.local.yml ps

# 后端健康检查
curl http://localhost:8080/health

# 数据库连接测试
docker exec sub2api-postgres psql -U sub2api -d sub2api -c "SELECT 1;"

# Redis 连接测试
docker exec sub2api-redis redis-cli ping
```

### 日志查看

```bash
# Docker 容器日志
docker compose -f docker-compose.local.yml logs -f postgres
docker compose -f docker-compose.local.yml logs -f redis

# 后端日志（直接输出到终端）
# 使用 ./server 或 go run 时日志直接显示
```

### 数据库操作

```bash
# 查看用户表
docker exec sub2api-postgres psql -U sub2api -d sub2api -c "SELECT id, email, role FROM users;"

# 重置管理员密码（需先生成 bcrypt hash）
# 生成 hash: go run /tmp/genhash.go
docker exec sub2api-postgres psql -U sub2api -d sub2api -c "UPDATE users SET password_hash='<new_hash>' WHERE email='admin@sub2api.local';"
```

---

## 七、故障排查

### 问题 1: 后端启动失败 - 密码认证错误

**错误信息:** `pq: password authentication failed for user "sub2api"`

**原因:** 
- `.env` 中缺少 `DATABASE_PASSWORD` 变量
- `config.yaml` 中密码配置错误

**解决方案:**
1. 确保 `.env` 中有 `DATABASE_PASSWORD=263fad824a6bae750f6ad5690cdd894f`
2. 确保 `config.yaml` 中 `database.password` 与 `.env` 中一致

### 问题 2: 管理员用户不存在

**原因:** 
- `config.yaml` 已存在，AUTO_SETUP 跳过管理员创建
- 数据库用户表为空

**解决方案:**
手动在数据库创建管理员用户（见 4.3 节）

### 问题 3: air 命令未找到

**解决方案:**
```bash
# 安装 air
go install github.com/air-verse/air@latest

# 添加到 PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## 八、文件路径汇总

| 文件 | 路径 | 用途 |
|------|------|------|
| 环境变量 | `/home/root/sub2api/deploy/.env` | Docker 容器和后端配置 |
| 后端配置 | `/home/root/sub2api/backend/config.yaml` | 后端服务配置 |
| Docker Compose | `/home/root/sub2api/deploy/docker-compose.local.yml` | 本地目录版部署配置 |
| Docker Compose (dev) | `/home/root/sub2api/deploy/docker-compose.dev-override.yml` | 开发模式端口暴露 |
| Dockerfile | `/home/root/sub2api/Dockerfile` | 多阶段构建镜像 |
| 构建脚本 | `/home/root/sub2api/deploy/build_image.sh` | 快速构建镜像脚本 |

---

*文档生成时间: 2026-06-13*