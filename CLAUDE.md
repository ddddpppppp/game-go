# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码库中工作时提供指导。

## 项目概述

这是一个基于 GoFrame v2.9.0 构建的 Go 应用程序，提供游戏 API 集成和客户对话管理功能。系统处理 Facebook 页面同步、游戏免费试玩/充值任务、基于 WebSocket 的实时通信，并与 SalesSmartly 等外部服务集成以提供客户服务。

## 开发命令

### 构建和运行
```bash
# 构建应用程序
make build
# 或手动构建：
go build -o ai-service-bot .

# 运行应用程序
go run main.go
# 服务器在 :8000 端口启动，Swagger UI 在 /swagger/
```

### 代码生成 (GoFrame CLI)
```bash
# 如果不存在则安装 GoFrame CLI
make cli.install

# 数据库模式更改后生成 DAO 文件
make dao
# 等同于：gf gen dao

# 从 API 定义生成控制器
make ctrl
# 等同于：gf gen ctrl

# 生成服务接口
make service
# 等同于：gf gen service

# 更新 GoFrame 到最新版本
make up
# 等同于：gf up -a
```

### 数据库操作
```bash
# 为特定数据库组生成 DAO 文件
gf gen dao -g ai  # 针对 ai_system 数据库
gf gen dao        # 针对默认数据库 (gf)
```

### Docker 和部署
```bash
# 构建 Docker 镜像
make image

# 构建并推送到镜像仓库
make image.push

# 部署到 Kubernetes
make deploy
```

## 架构概述

### 框架和依赖
- **GoFrame v2.9.0**：主要的 Web 框架，提供 HTTP 服务器、ORM、验证、日志、缓存
- **MySQL**：两个数据库 - `gf`（主库）和 `ai_system`（AI 相关数据）
- **Redis**：缓存和定时任务的分布式锁
- **WebSocket**：实时通信和连接管理
- **RabbitMQ**：消息队列（代码中当前已禁用）

### 目录结构
```
internal/
├── cmd/                    # 应用程序入口点和命令定义
├── controller/             # HTTP 请求处理器（game_api、game_conversation、game_conversation_ws）
├── service/               # 业务逻辑层，包含 Redis 缓存
├── dao/                   # 自动生成的数据访问对象
├── model/                 # 实体、DO（数据对象）、数据库模型
├── crontab/               # 使用 Redis 锁的后台任务处理
├── middleware/            # 自定义响应处理器、认证、追踪
└── consts/               # 应用程序常量和 Redis 键

api/                       # 带路由标签的 RESTful API 定义
boot/                      # 应用程序初始化（日志、OSS、RabbitMQ）
manifest/config/           # YAML 配置文件
```

### 关键模式

#### 控制器模式
控制器处理 HTTP 请求并委托给服务：
```go
func (c *ControllerV1) GameFreeplay(ctx context.Context, req *v1.GameFreeplayReq) (res *v1.GameCommonRes, err error)
```

#### 带缓存的服务层
服务实现业务逻辑并使用 Redis 缓存：
- 使用 `g.Redis()` 进行缓存，设置适当的 TTL
- 缓存键定义在 `internal/consts/` 中
- 遵循缓存旁路模式进行数据访问

#### 后台任务处理
定时任务使用基于 Redis 的分布式锁：
- 通过 `internal/cmd/cron.go` 中的 `gcron.Add()` 调度
- 每种任务类型在 `internal/crontab/` 中都有自己的定时服务
- Redis 锁防止跨实例的重复处理

#### WebSocket 连接管理
- 每个用户支持多连接，具有自动清理功能
- 使用互斥锁进行线程安全操作
- WebSocket 连接需要管理员认证

### 数据库架构

#### 多数据库支持
- **默认组**（`gf` 数据库）：主应用程序数据
- **AI 组**（`ai_system` 数据库）：AI 相关数据
- 在 `manifest/config/config.yaml` 的 `database` 部分配置

#### DAO 生成
- 模式更改后运行 `make dao`
- 自动生成 Entity、DO 和 DAO 文件
- 支持使用 `-g` 标志的多数据库组

### 配置管理

#### 配置文件
- 主配置：`manifest/config/config.yaml`
- 开发配置生成：`hack/config.yaml`
- 环境特定覆盖：`manifest/deploy/kustomize/overlays/`

#### 关键配置部分
- `server`：HTTP 服务器设置（端口 8000）
- `database`：多数据库连接
- `redis`：缓存配置
- `salessmartly`：外部服务集成
- `oss`：阿里云对象存储服务
- `game.api.mock`：开发模式启用模拟

### 外部服务集成

#### SalesSmartly
客户服务平台集成，支持 webhook：
- API 主机、webhook 主机和认证在 YAML 中配置
- 服务实现在 `internal/service/sales_smartly/`

#### 游戏 API 集成
- 通过 `game.api.mock` 配置提供模拟模式
- 服务层在 `internal/service/game_api/`

### 开发指南

#### 添加新功能
1. 在 `api/` 目录中定义 API 结构
2. 运行 `make ctrl` 生成控制器样板
3. 在服务层实现业务逻辑
4. 在适当的地方使用 Redis 添加缓存
5. 如果涉及外部服务，更新配置

#### 数据库更改
1. 更新数据库模式
2. 运行 `make dao` 重新生成 DAO 文件
3. 更新服务中的模型引用

#### 后台任务
1. 在 `internal/crontab/` 中创建定时服务
2. 在 `internal/cmd/cron.go` 中添加定时调度
3. 实现基于 Redis 的分布式处理锁

#### WebSocket 功能
1. 使用 `internal/service/game_conversation_ws/` 中现有的连接管理器
2. 实现消息广播模式
3. 确保正确的连接清理

### 测试和调试

#### API 测试
- Swagger UI 可在 `http://localhost:8000/swagger/` 访问
- OpenAPI 规范在 `http://localhost:8000/api.json`

#### 数据库调试
- 在配置中设置 `database.default.debug: true` 以启用 SQL 查询日志
- 日志显示在控制台和日志文件中

#### Redis 监控
- 监控缓存和锁定行为的 Redis 键
- 定时任务状态通过带 TTL 的 Redis 锁跟踪

### 常见问题和解决方案

#### CLI 工具安装
如果找不到 `gf` 命令，运行 `make cli.install` 自动安装 GoFrame CLI。

#### 数据库连接问题
检查 `manifest/config/config.yaml` 中的正确数据库凭据，确保 MySQL 正在运行。

#### 定时任务冲突
后台任务使用基于 Redis 的锁定。如果任务未运行，检查 Redis 连接和锁过期时间。

#### WebSocket 连接问题
确保管理员认证中间件配置正确，Redis 可用于会话管理。