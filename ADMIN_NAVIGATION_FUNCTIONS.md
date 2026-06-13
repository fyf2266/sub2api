# Sub2API 管理后台导航功能操作说明

## 文档概述
本文档详细描述 Sub2API 管理后台左侧导航栏中每个功能模块的具体使用方法、操作流程及其关联功能。

---

## 一、管理区（Admin Section）

### 1. 仪表盘（Dashboard）
**路径**: `/admin/dashboard`

**功能定位**: 系统总览与数据统计

**具体使用方法**:
- 查看系统整体运行状态和关键指标
- 实时监控 API 请求量、成功率、延迟等核心指标
- 查看用户增长趋势、使用量统计
- 监控系统资源使用情况

**关联功能**:
- **运维监控**: 提供更详细的实时监控数据
- **使用记录**: 可查看详细的请求日志
- **系统设置**: 配置仪表盘显示参数

**后端实现**: [dashboard_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/dashboard_handler.go)

---

### 2. 运维监控（Ops Monitor）
**路径**: `/admin/ops`

**功能定位**: 系统运维与实时监控

**具体使用方法**:
- 实时查看 QPS/TPS 数据（WebSocket 推送）
- 查看错误日志和请求错误详情
- 监控上游服务状态
- 配置高级设置和指标阈值
- 查看系统日志和健康状态

**关联功能**:
- **仪表盘**: 汇总展示监控数据
- **渠道监控**: 监控各渠道运行状态
- **系统设置**: 配置监控参数

**后端实现**: [ops_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/ops_handler.go)、[ops_dashboard_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/ops_dashboard_handler.go)

---

### 3. 用户管理（Users）
**路径**: `/admin/users`

**功能定位**: 管理平台用户账户

**具体使用方法**:
- 创建新用户（邮箱、密码、用户名）
- 编辑用户信息（余额、并发数、RPM限制）
- 设置用户允许访问的分组
- 配置用户专属分组倍率
- 管理用户状态（启用/禁用）
- 绑定用户认证身份（OAuth）

**关联功能**:
- **分组管理**: 用户与分组关联
- **订阅管理**: 用户订阅分配
- **API密钥**: 用户的 API Key 管理
- **账户管理**: 用户账户统计

**后端实现**: [user_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/user_handler.go)

---

### 4. 分组管理（Groups）
**路径**: `/admin/groups`

**功能定位**: 管理 AI 服务分组配置

**具体使用方法**:
- 创建分组（名称、平台、描述）
- 配置分组倍率（Rate Multiplier）
- 设置分组类型（标准/订阅）
- 配置消费限额（日/周/月）
- 管理图片生成计费（Antigravity/Gemini）
- 批量设置分组倍率和 RPM 覆盖
- 查看分组统计和 API Key

**关联功能**:
- **渠道管理**: 分组与渠道关联
- **用户管理**: 用户与分组关联
- **订阅管理**: 订阅分组配置
- **账号管理**: 账号与分组关联

**后端实现**: [group_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/group_handler.go)

---

### 5. 渠道管理（Channels）

#### 5.1 渠道定价（Channel Pricing）
**路径**: `/admin/channels/pricing`

**功能定位**: 配置渠道级别的定价策略

**具体使用方法**:
- 创建渠道（名称、描述、关联分组）
- 配置模型定价（输入/输出价格、缓存读写价格）
- 设置模型映射关系
- 配置计费模式（token/请求/图片）
- 设置分级定价规则

**关联功能**:
- **分组管理**: 渠道关联分组
- **账号管理**: 渠道关联账号
- **渠道监控**: 监控渠道状态

#### 5.2 渠道监控（Channel Monitor）
**路径**: `/admin/channels/monitor`

**功能定位**: 监控各渠道的运行状态

**具体使用方法**:
- 实时查看各渠道请求量和成功率
- 监控渠道延迟和错误率
- 配置监控告警规则
- 查看渠道请求模板

**关联功能**:
- **运维监控**: 全局监控
- **账号管理**: 账号状态监控

**后端实现**: [channel_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/channel_handler.go)、[channel_monitor_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/channel_monitor_handler.go)

---

### 6. 订阅管理（Subscriptions）
**路径**: `/admin/subscriptions`

**功能定位**: 管理用户订阅服务

**具体使用方法**:
- 查看用户订阅列表
- 分配订阅给用户
- 批量分配订阅
- 延长订阅有效期
- 查看订阅进度和状态

**关联功能**:
- **用户管理**: 用户订阅分配
- **分组管理**: 订阅分组配置
- **账号管理**: 订阅账号关联

**后端实现**: [subscription_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/subscription_handler.go)

---

### 7. 账号管理（Accounts）
**路径**: `/admin/accounts`

**功能定位**: 管理上游 AI 平台账号

**具体使用方法**:
- 创建上游账号（API Key/OAuth）
- 配置账号并发数和限额
- 管理账号状态（启用/禁用/调度）
- 导入 Codex Session
- 同步账号信息
- 检查混合渠道配置

**关联功能**:
- **分组管理**: 账号与分组关联
- **渠道管理**: 账号与渠道关联
- **运维监控**: 账号状态监控

**后端实现**: [account_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/account_handler.go)

---

### 8. 公告（Announcements）
**路径**: `/admin/announcements`

**功能定位**: 管理系统公告

**具体使用方法**:
- 创建公告（标题、内容、优先级）
- 编辑和删除公告
- 设置公告生效时间
- 管理公告状态

**关联功能**:
- 无直接关联功能，独立模块

**后端实现**: [announcement_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/announcement_handler.go)

---

### 9. IP管理（Proxies）
**路径**: `/admin/proxies`

**功能定位**: 管理代理 IP 配置

**具体使用方法**:
- 添加代理服务器配置
- 配置代理认证信息
- 管理代理状态
- 配置 IP 白名单/黑名单

**关联功能**:
- **账号管理**: 账号代理配置
- **渠道管理**: 渠道代理配置

**后端实现**: [proxy_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/proxy_handler.go)

---

### 10. 兑换码（Redeem）
**路径**: `/admin/redeem`

**功能定位**: 管理充值兑换码

**具体使用方法**:
- 生成兑换码（面值、数量）
- 配置兑换码有效期
- 查看兑换码使用状态
- 导出兑换码列表

**关联功能**:
- **用户管理**: 用户兑换记录
- **订阅管理**: 订阅兑换

**后端实现**: [redeem_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/redeem_handler.go)

---

### 11. 优惠码（Promo Codes）
**路径**: `/admin/promo-codes`

**功能定位**: 管理促销优惠码

**具体使用方法**:
- 创建优惠码（折扣比例、使用次数）
- 配置优惠码适用范围
- 设置有效期
- 查看使用统计

**关联功能**:
- **用户管理**: 用户使用记录
- **订阅管理**: 订阅优惠

**后端实现**: [promo_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/promo_handler.go)

---

### 12. 使用记录（Usage）
**路径**: `/admin/usage`

**功能定位**: 查看系统使用记录

**具体使用方法**:
- 查询用户使用记录
- 按时间范围筛选
- 按用户/分组/模型筛选
- 导出使用报表

**关联功能**:
- **仪表盘**: 汇总统计
- **用户管理**: 用户详细记录
- **分组管理**: 分组使用统计

**后端实现**: [usage_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/usage_handler.go)

---

### 13. 系统设置（Settings）
**路径**: `/admin/settings`

**功能定位**: 系统全局配置

**具体使用方法**:
- 配置平台基础设置
- 设置认证源默认值
- 配置平台配额
- 配置钉钉通知
- 管理系统参数

**关联功能**:
- **运维监控**: 监控参数配置
- **渠道管理**: 渠道默认配置

**后端实现**: [setting_handler.go](file:///home/root/sub2api/backend/internal/handler/admin/setting_handler.go)

---

## 二、个人账户区（Personal Section）

### 1. API 密钥（API Keys）
**路径**: `/keys`

**功能定位**: 管理个人 API Key

**具体使用方法**:
- 生成新的 API Key
- 配置 Key 权限和限额
- 管理 Key 状态（启用/禁用）
- 查看 Key 使用统计

**关联功能**:
- **用户管理**: 管理员视角的 Key 管理
- **使用记录**: Key 的使用日志

---

### 2. 使用记录（Usage）
**路径**: `/usage`

**功能定位**: 查看个人使用记录

**具体使用方法**:
- 查看个人 API 使用统计
- 按时间范围查询
- 查看消费明细

**关联功能**:
- **管理员使用记录**: 管理员视角的全局记录

---

### 3. 渠道状态（Monitor）
**路径**: `/monitor`

**功能定位**: 查看个人可用渠道状态

**具体使用方法**:
- 查看各渠道健康状态
- 了解可用模型和定价
- 监控个人请求状态

**关联功能**:
- **运维监控**: 管理员视角的渠道监控

---

### 4. 我的订阅（Subscriptions）
**路径**: `/subscriptions`

**功能定位**: 管理个人订阅

**具体使用方法**:
- 查看当前订阅状态
- 查看订阅有效期
- 管理订阅配置

**关联功能**:
- **订阅管理**: 管理员视角的订阅管理

---

### 5. 兑换（Redeem）
**路径**: `/redeem`

**功能定位**: 使用兑换码充值

**具体使用方法**:
- 输入兑换码进行充值
- 查看兑换历史记录

**关联功能**:
- **兑换码管理**: 管理员视角的兑换码管理

---

### 6. 个人资料（Profile）
**路径**: `/profile`

**功能定位**: 管理个人信息

**具体使用方法**:
- 修改个人资料（用户名、邮箱）
- 更新密码
- 管理通知偏好

**关联功能**:
- **用户管理**: 管理员视角的用户信息管理

---

## 三、功能关系图

```
┌─────────────────────────────────────────────────────────────────┐
│                        核心业务流程                              │
├─────────────────────────────────────────────────────────────────┤
│  账号管理 ←→ 分组管理 ←→ 渠道管理 ←→ 用户管理                     │
│     ↑           ↑           ↑           ↑                       │
│     │           │           │           │                       │
│  运维监控 ←─── 仪表盘 ───→ 使用记录 ←─── 订阅管理                  │
│     │                                      │                     │
│     ↓                                      ↓                     │
│  渠道监控                              兑换码/优惠码               │
└─────────────────────────────────────────────────────────────────┘
```

---

## 四、权限说明

| 功能模块 | 管理员权限 | 普通用户权限 |
|---------|-----------|-------------|
| 仪表盘 | 完整访问 | 只读访问 |
| 运维监控 | 完整访问 | 无访问 |
| 用户管理 | 完整访问 | 无访问 |
| 分组管理 | 完整访问 | 无访问 |
| 渠道管理 | 完整访问 | 无访问 |
| 订阅管理 | 完整访问 | 仅查看个人 |
| 账号管理 | 完整访问 | 无访问 |
| 公告 | 完整访问 | 只读访问 |
| IP管理 | 完整访问 | 无访问 |
| 兑换码 | 完整访问 | 仅兑换 |
| 优惠码 | 完整访问 | 仅使用 |
| 使用记录 | 完整访问 | 仅查看个人 |
| 系统设置 | 完整访问 | 无访问 |

---

## 文档版本

- **版本**: v1.0
- **生成时间**: 2026-06-13
- **适用系统**: Sub2API Admin Dashboard
