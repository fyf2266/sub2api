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

## 八、生产环境部署（自建 Docker 镜像）

> 本节介绍如何根据本地代码重新编译 Docker 镜像并部署生产环境

### 8.1 前置条件

- Docker 20.10+
- Docker Compose v2+
- Git（用于获取代码）

### 8.2 部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                      生产环境架构                              │
│                                                             │
│   Docker 网络 (sub2api-network):                            │
│   ├── sub2api (自建镜像 :production)                       │
│   │   端口: 8080 (仅内部访问，由 Nginx 代理)                │
│   │                                                   │
│   ├── PostgreSQL 18                                        │
│   │   端口: 5432 (仅内部访问)                             │
│   │   数据目录: ./postgres_data                          │
│   │                                                       │
│   └── Redis 8                                             │
│       端口: 6379 (仅内部访问)                              │
│       数据目录: ./redis_data                              │
│                                                             │
│   外部访问: 通过 Nginx 反向代理到 sub2api:8080              │
└─────────────────────────────────────────────────────────────┘
```

### 8.3 构建 Docker 镜像

**步骤 1：获取源码**

```bash
cd /home/root/sub2api

# 如果使用现有代码，直接进入下一步
# 如果需要拉取最新代码：
git pull origin main
```

**步骤 2：构建生产镜像**

```bash
# 在项目根目录执行构建
docker build -t sub2api:production \
  --build-arg VERSION=$(git rev-parse --short HEAD 2>/dev/null || echo "dev") \
  --build-arg DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  .

# 或者构建 latest 标签
docker build -t sub2api:latest .
```

**步骤 3：验证镜像构建**

```bash
# 查看构建的镜像
docker images sub2api

# 运行测试容器验证
docker run --rm sub2api:production /app/sub2api version
```

### 8.4 配置生产环境

**步骤 1：创建生产环境目录**

```bash
mkdir -p /opt/sub2api-production && cd /opt/sub2api-production
```

**步骤 2：复制必要文件**

```bash
# 复制 docker-compose 配置
cp /home/root/sub2api/deploy/docker-compose.local.yml .
cp /home/root/sub2api/deploy/docker-compose.production.yml .

# 复制环境变量模板
cp /home/root/sub2api/deploy/.env.example .env
```

**步骤 3：配置环境变量**

```bash
# 生成安全密钥
JWT_SECRET=$(openssl rand -hex 32)
TOTP_KEY=$(openssl rand -hex 32)
POSTGRES_PASSWORD=$(openssl rand -hex 32)

# 编辑 .env 文件
cat > .env << EOF
# =============================================================================
# Sub2API 生产环境配置
# =============================================================================

# 服务器配置
SERVER_MODE=release
LOG_LEVEL=info
LOG_FORMAT=json
LOG_ENV=production
TZ=Asia/Shanghai

# 数据库配置
POSTGRES_USER=sub2api
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
DATABASE_PASSWORD=${POSTGRES_PASSWORD}
POSTGRES_DB=sub2api
DATABASE_MAX_OPEN_CONNS=256
DATABASE_MAX_IDLE_CONNS=128

# Redis 配置
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=4096

# JWT 配置
JWT_SECRET=${JWT_SECRET}
JWT_EXPIRE_HOUR=24

# TOTP 配置
TOTP_ENCRYPTION_KEY=${TOTP_KEY}

# 安全配置
SECURITY_URL_ALLOWLIST_ENABLED=false
SECURITY_URL_ALLOWLIST_ALLOW_INSECURE_HTTP=false
SECURITY_URL_ALLOWLIST_ALLOW_PRIVATE_HOSTS=false
EOF
```

**步骤 4：创建数据目录**

```bash
mkdir -p data postgres_data redis_data
chmod -R 755 .
chmod 600 .env
```

### 8.5 启动服务

```bash
# 启动所有服务（生产模式）
docker compose -f docker-compose.local.yml -f docker-compose.production.yml up -d

# 查看服务状态
docker compose -f docker-compose.local.yml ps

# 查看日志
docker compose -f docker-compose.local.yml logs -f sub2api
```

**验证服务启动：**

```bash
# 健康检查
curl http://localhost:8080/health

# 查看容器状态
docker compose -f docker-compose.local.yml ps
```

### 8.6 Nginx 反向代理配置（推荐）

```nginx
# /etc/nginx/conf.d/sub2api.conf

upstream sub2api_backend {
    server 127.0.0.1:8080;
    keepalive 64;
}

server {
    listen 80;
    server_name your-domain.com;

    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/fullchain.pem;
    ssl_certificate_key /path/to/privkey.pem;
    ssl_session_timeout 1d;
    ssl_session_cache shared:SSL:50m;
    ssl_session_tickets off;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    underscores_in_headers on;

    client_max_body_size 256m;

    location / {
        proxy_pass http://sub2api_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Connection "";

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        proxy_buffering off;
        proxy_cache off;
    }
}
```

### 8.7 生产环境运维

#### 8.7.1 启动/停止/重启

```bash
# 启动服务
docker compose -f docker-compose.local.yml -f docker-compose.production.yml up -d

# 停止服务
docker compose -f docker-compose.local.yml -f docker-compose.production.yml down

# 重启服务
docker compose -f docker-compose.local.yml -f docker-compose.production.yml restart

# 查看状态
docker compose -f docker-compose.local.yml -f docker-compose.production.yml ps
```

#### 8.7.2 更新服务

当代码有更新时，重新构建镜像并重启：

```bash
# 1. 进入代码目录
cd /home/root/sub2api

# 2. 拉取最新代码（如果有）
git pull origin main

# 3. 重新构建镜像
docker build -t sub2api:production \
  --build-arg VERSION=$(git rev-parse --short HEAD) \
  --build-arg DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  .

# 4. 重启服务
docker compose -f docker-compose.local.yml -f docker-compose.production.yml up -d
```

#### 8.7.3 日志管理

```bash
# 实时查看所有日志
docker compose -f docker-compose.local.yml logs -f

# 实时查看 sub2api 日志
docker compose -f docker-compose.local.yml logs -f sub2api

# 查看最近 100 行
docker compose -f docker-compose.local.yml logs --tail 100

# 查看特定容器日志
docker compose -f docker-compose.local.yml logs -f postgres
docker compose -f docker-compose.local.yml logs -f redis
```

#### 8.7.4 数据备份

```bash
# 备份所有数据
tar czf sub2api-backup-$(date +%Y%m%d-%H%M%S).tar.gz data/ postgres_data/ redis_data/ .env

# 备份数据库（SQL 格式）
docker exec sub2api-postgres pg_dump -U sub2api sub2api > database-backup-$(date +%Y%m%d).sql
```

#### 8.7.5 数据恢复

```bash
# 停止服务
docker compose -f docker-compose.local.yml -f docker-compose.production.yml down

# 解压备份
tar xzf sub2api-backup-20250613-120000.tar.gz

# 恢复数据库
cat database-backup-20250613.sql | docker exec -i sub2api-postgres psql -U sub2api sub2api

# 重启服务
docker compose -f docker-compose.local.yml -f docker-compose.production.yml up -d
```

---

## 九、文件路径汇总

| 文件 | 路径 | 用途 |
|------|------|------|
| 环境变量 | `/home/root/sub2api/deploy/.env` | Docker 容器和后端配置 |
| 后端配置 | `/home/root/sub2api/backend/config.yaml` | 后端服务配置 |
| Docker Compose | `/home/root/sub2api/deploy/docker-compose.local.yml` | 本地目录版部署配置 |
| Docker Compose (dev) | `/home/root/sub2api/deploy/docker-compose.dev-override.yml` | 开发模式端口暴露 |
| Docker Compose (prod) | `/home/root/sub2api/deploy/docker-compose.production.yml` | 生产模式配置 |
| Dockerfile | `/home/root/sub2api/Dockerfile` | 多阶段构建镜像 |
| 构建脚本 | `/home/root/sub2api/deploy/build_image.sh` | 快速构建镜像脚本 |

---

*文档生成时间: 2026-06-13*