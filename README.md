# gin-layout

基于 **Gin 框架**的标准分层 Go 后端项目模板，提供开箱即用的 RBAC 权限管理、JWT 双令牌认证、多数据库 ORM 支持、结构化日志、以及完整的接口文档生成。

## 📋 核心特性

### 架构设计

- **分层架构** — Handler → Service → Repository 三层模式，清晰的职责分离
- **依赖注入** — bootstrap 模块手动管理依赖，无第三方 DI 框架，依赖关系一目了然
- **泛型 CRUD** — 使用 Go 1.18+ 泛型实现 `CRUDRepository[T, ID]`，避免重复代码

### 身份认证与权限

- **JWT 双令牌认证** — Access Token (5分钟) + Refresh Token (12小时)，支持令牌撤销黑名单
- **RBAC 权限管理** — 基于 Casbin v3，支持角色、权限、菜单的灵活配置，GORM 持久化
- **默认超级管理员** — 首次启动自动创建超管角色、账号 (admin/admin123) 和权限

### 数据访问

- **多数据库支持** — 开箱即用支持 MySQL、PostgreSQL、SQLite、SQL Server
- **自动表迁移** — 启动时执行 GORM AutoMigrate，无需手动 SQL 脚本
- **类型安全的查询** — 完整的错误转换和类型检查

### 错误处理与日志

- **业务错误体系** — 自定义 `DomainError` 类型（含业务码、HTTP 状态、消息），支持 `errors.Is` 判定
- **统一响应格式** — 所有 HTTP 响应统一为 `{code, message, data}` 结构
- **结构化日志** — 基于 zerolog，支持控制台/JSON 格式、灵活的日志级别控制
- **请求日志中间件** — 自动记录请求/响应信息，支持请求 ID 链路追踪

### 安全与操作

- **密码加密** — 采用 bcrypt 哈希存储，安全等级高
- **CORS 中间件** — 可配置的跨域资源共享
- **优雅关闭** — 支持 SIGINT/SIGTERM 信号处理，安全关闭 HTTP 服务
- **Docker 部署** — 完整的 Dockerfile (多阶段构建) + docker-compose 配置

### 文档与接口

- **API 文档 (Redoc)** — 使用 Redoc 提供交互式 API 文档
- **自动化文档生成** — 支持自定义 Schema 和接口分组

## 🛠 技术栈

| 分类           | 技术                               | 版本  | 用途           |
| -------------- | ---------------------------------- | ----- | -------------- |
| **Web 框架**   | Gin                                | v1.10 | HTTP 服务器    |
| **ORM**        | GORM                               | v1.31 | 数据库操作     |
| **数据库驱动** | MySQL/PostgreSQL/SQLite/SQL Server | -     | 多数据库支持   |
| **权限控制**   | Casbin                             | v3    | RBAC 权限模型  |
| **JWT 认证**   | golang-jwt/jwt                     | v5    | 令牌生成和验证 |
| **缓存**       | Redis                              | v9    | 会话和黑名单   |
| **日志**       | rs/zerolog                         | v1.33 | 结构化日志     |
| **配置管理**   | spf13/viper                        | v1.21 | TOML 配置加载  |
| **密码加密**   | golang.org/x/crypto                | -     | bcrypt 算法    |
| **接口文档**   | swaggo/swag                        | v1.16 | API 文档生成   |
| **工具库**     | samber/lo                          | -     | 函数式工具集   |
| **类型转换**   | spf13/cast                         | v1.10 | 类型安全转换   |

**Go 版本**：1.26.3 或更高

## 📂 项目结构

```
gin-layout/
├── cmd/
│   └── server/
│       └── main.go                      # 应用入口：加载配置 → 初始化应用 → 启动服务 → 优雅关闭
│
├── config/                              # 配置加载与定义
│   ├── conf.go                          # Viper 从 TOML 加载配置
│   ├── database.go                      # 数据库 DSN 构建器
│   └── types.go                         # 配置结构体定义
│
├── etc/                                 # 环境配置文件
│   ├── config.toml                      # 本地开发配置
│   ├── config.docker.toml               # Docker 环境配置
│   └── casbin_model.conf                # Casbin RBAC 权限模型
│
├── internal/                            # 核心业务代码（禁止外部导入）
│   ├── bootstrap/                       # 应用启动与依赖注入
│   │   ├── app.go                       # 应用初始化：logger → DB → migrate → seed → security → repos → services → handlers → router → server
│   │   ├── initializer/
│   │   │   ├── initializer.go           # 初始化管理器
│   │   │   └── menu.json                # 默认菜单配置
│   │   └── seed/
│   │       └── admin.go                 # 创建默认超管账号和角色
│   │
│   ├── auth/                            # 认证模块
│   │   ├── handler.go                   # HTTP 处理器：登录、刷新令牌、登出
│   │   ├── service.go                   # 认证业务逻辑
│   │   ├── dto.go                       # 数据传输对象
│   │   └── ports.go                     # 依赖接口定义
│   │
│   ├── user/                            # 用户管理
│   │   ├── handler.go                   # 用户 CRUD HTTP 接口
│   │   ├── service.go                   # 用户业务逻辑
│   │   ├── repo.go                      # 用户数据访问
│   │   ├── dto.go                       # DTO 定义
│   │   └── ports.go                     # 接口定义
│   │
│   ├── role/                            # 角色管理
│   │   ├── handler.go                   # 角色 CRUD HTTP 接口
│   │   ├── service.go                   # 角色业务逻辑
│   │   ├── repo.go                      # 角色数据访问
│   │   ├── dto.go                       # DTO 定义
│   │   └── ports.go                     # 接口定义
│   │
│   ├── menu/                            # 菜单管理（树形结构）
│   │   ├── handler.go                   # 菜单 CRUD HTTP 接口
│   │   ├── service.go                   # 菜单业务逻辑
│   │   ├── repo.go                      # 菜单数据访问
│   │   ├── helper.go                    # 树形结构处理辅助函数
│   │   ├── dto.go                       # DTO 定义
│   │   └── ports.go                     # 接口定义
│   │
│   ├── health/                          # 健康检查
│   │   ├── handler.go                   # /health 接口
│   │   ├── service.go                   # 健康检查逻辑
│   │   └── ports.go                     # 接口定义
│   │
│   ├── middleware/                      # Gin 中间件
│   │   ├── auth.go                      # JWT 认证中间件
│   │   ├── auth_test.go                 # 认证中间件单元测试
│   │   ├── rbac.go                      # Casbin 权限检查中间件
│   │   ├── cors.go                      # CORS 跨域中间件
│   │   ├── cors_test.go                 # CORS 测试
│   │   ├── error.go                     # 错误响应中间件
│   │   ├── logger.go                    # 请求日志中间件
│   │   ├── recovery.go                  # 恐慌恢复中间件
│   │   └── request_id.go                # 请求 ID 链路追踪中间件
│   │
│   ├── domain/                          # 领域模型和错误定义
│   │   ├── user.go                      # User 聚合根
│   │   ├── role.go                      # Role 聚合根
│   │   ├── menu.go                      # Menu 聚合根
│   │   ├── join.go                      # 多对多关联表定义
│   │   ├── page.go                      # 分页请求/响应（泛型）
│   │   ├── context.go                   # CurrentUser 上下文传递
│   │   ├── error.go                     # DomainError 业务错误定义
│   │   ├── errors.go                    # 业务错误常量
│   │   └── token.go                     # Token 相关定义
│   │
│   ├── infra/                           # 基础设施层
│   │   ├── database.go                  # 数据库连接初始化
│   │   ├── crud.go                      # 泛型 CRUDRepository 实现
│   │   ├── jwt.go                       # JWT 令牌服务
│   │   ├── logger.go                    # zerolog 日志初始化
│   │   ├── password.go                  # bcrypt 密码加密/验证
│   │   ├── policy.go                    # Casbin 权限策略管理
│   │   ├── cache.go                     # Redis 缓存
│   │   ├── token_blacklist.go           # 令牌黑名单实现
│   │   └── errors.go                    # GORM 错误转换
│   │
│   ├── router/                          # 路由配置
│   │   ├── router.go                    # Router 接口定义
│   │   └── admin.go                     # 管理后台路由组
│   │
│   ├── server/                          # HTTP 服务器
│   │   ├── server.go                    # HTTP 服务启动、优雅关闭
│   │   └── endpoint.go                  # TypedHandler 泛型包装器
│   │
│   ├── apidoc/                          # API 文档生成模块
│   │   ├── builder.go                   # 文档构建器
│   │   ├── config.go                    # 文档配置
│   │   ├── model.go                     # 文档数据模型
│   │   ├── publisher.go                 # 文档发布器
│   │   ├── registry.go                  # Schema 注册表
│   │   ├── schema_builder.go            # Schema 构建器
│   │   ├── swagger2_renderer.go         # Swagger 2.0 渲染
│   │   └── testdata/                    # 测试数据
│   │
│   └── common/                          # 通用工具
│       ├── response.go                  # 统一响应格式
│       ├── base.go                      # 基础模型
│       ├── interfaces.go                # 通用接口定义
│       └── token.go                     # Token 工具函数
│
├── deploy/                              # 部署配置
│   ├── Dockerfile                       # 多阶段构建
│   └── docker-compose.yml               # 完整栈：应用 + MySQL + Redis
│
├── web/                                 # 前端项目（Vue 3 + Naive UI）
│   ├── package.json                     # 前端依赖
│   ├── vite.config.js                   # Vite 构建配置
│   ├── tsconfig.json                    # TypeScript 配置
│   ├── uno.config.js                    # UnoCSS 原子化 CSS 配置
│   ├── src/                             # 前端源代码
│   └── public/                          # 公开资源
│
├── docs/                                # 文档和规划
│   ├── superpowers/
│   │   ├── plans/                       # 功能规划
│   │   └── specs/                       # 技术方案
│   └── README.md                        # 文档索引
│
├── go.mod                               # Go 模块定义
├── Makefile                             # 构建命令
└── README.md                            # 本文件
```

## 🚀 快速开始

### 环境要求

- Go 1.26.3 或更高
- MySQL 5.7+ 或 PostgreSQL 12+ 或 SQLite 3 或 SQL Server 2017+
- Redis 6.0+（可选，用于缓存和令牌黑名单）
- Node.js 18+ 和 pnpm（开发前端时）

### 本地开发

#### 1. 克隆项目

```bash
git clone <repository-url>
cd gin-layout
```

#### 2. 配置数据库

编辑 `etc/config.toml`，配置数据库连接：

```toml
[server]
host = "0.0.0.0"
port = 8085
mode = "debug"

[database]
driver = "mysql"
host = "127.0.0.1"
port = 3306
user = "root"
password = "root"
database = "gin_layout"

[redis]
mode = "single"
addrs = ["127.0.0.1:6379"]
```

#### 3. 安装依赖

```bash
go mod download
go mod tidy
```

#### 4. 运行服务

```bash
# 开发模式
make dev

# 或编译后运行
make build
./bin/gin-layout -c etc/config.toml
```

#### 5. 访问服务

- **API 文档**：http://localhost:8085/swagger/index.html
- **健康检查**：http://localhost:8085/health
- **默认登录**：admin / admin123

### Docker 部署

```bash
# 一键启动完整栈（应用 + MySQL + Redis）
make docker-build-up

# 查看日志
docker-compose -f deploy/docker-compose.yml logs -f app

# 停止所有服务
make docker-down
```

## 💡 使用示例

### 登录与获取令牌

```bash
POST /login
Content-Type: application/json

{
  "account": "admin",
  "password": "admin123"
}

Response:
{
  "code": 0,
  "message": "ok",
  "data": {
    "accessToken": "eyJ...",
    "refreshToken": "eyJ...",
    "user": {...}
  }
}
```

### 调用受保护的 API

```bash
GET /api/v1/users?page=1&pageSize=10
Authorization: Bearer <accessToken>

Response:
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

### 刷新令牌

```bash
POST /refresh
Content-Type: application/json

{
  "refreshToken": "<refreshToken>"
}

Response:
{
  "code": 0,
  "message": "ok",
  "data": {
    "accessToken": "eyJ...",
    "refreshToken": "eyJ..."
  }
}
```

## 🧪 测试

```bash
# 运行所有测试
make test

# 生成覆盖率报告
make test-cover

# 代码检查
make lint
make vet
```

## 📊 API 响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

### 错误响应

```json
{
  "code": 40001,
  "message": "业务错误提示信息"
}
```

## 🔐 默认用户

| 账号  | 密码     | 角色       |
| ----- | -------- | ---------- |
| admin | admin123 | 超级管理员 |

## 🔑 配置说明

### JWT 配置 (`etc/config.toml`)

```toml
[jwt]
secret = "your-secret-key-change-in-production"
accessExpired = 300        # Access Token 有效期（秒）
refreshExpired = 43200     # Refresh Token 有效期（秒）
```

### 日志配置

```toml
[log]
level = "debug"            # debug, info, warn, error
format = "console"         # console 或 json
outputPath = "logs"        # 日志输出目录
```

### CORS 配置

```toml
[cors]
allowOrigins = ["http://localhost:5173"]
allowMethods = ["GET", "POST", "PUT", "DELETE"]
allowHeaders = ["Content-Type", "Authorization"]
allowCredentials = true
```

## 📖 Makefile 命令

```bash
make build          # 编译项目
make run            # 编译并运行
make dev            # 开发模式（直接运行）
make test           # 运行测试
make test-cover     # 生成覆盖率报告
make clean          # 清理构建产物
make lint           # 运行 golangci-lint
make vet            # 运行 go vet
make docker         # 构建 Docker 镜像
make docker-up      # 启动 Docker Compose
make docker-down    # 停止 Docker Compose
make docker-build-up # 构建镜像并启动
```

## 🏗 架构说明

### 分层架构

```
HTTP Handler
    ↓
Service (业务逻辑)
    ↓
Repository (数据访问)
    ↓
Database (GORM)
```

### 中间件顺序

1. RequestID - 生成请求 ID
2. Logger - 记录日志
3. CORS - 跨域处理
4. Auth - JWT 认证
5. RBAC - 权限检查
6. Recovery - 恐慌恢复
7. Error - 错误响应

### 启动顺序

1. 加载配置
2. 初始化日志
3. 连接数据库
4. 自动迁移表结构
5. 创建种子数据（默认超管）
6. 初始化安全组件（JWT、Casbin、Redis）
7. 创建数据访问层
8. 创建业务层
9. 创建 HTTP 处理器
10. 注册路由
11. 启动 HTTP 服务

## 🐛 常见问题

**Q: 如何切换数据库？**
A: 修改 `etc/config.toml` 中的 `[database]` 配置，支持 mysql、postgresql、sqlite、sqlserver

**Q: JWT 令牌过期了怎么办？**
A: 使用 Refresh Token 调用 `/refresh` 端点获取新的 Access Token

**Q: 如何自定义业务错误码？**
A: 在 `internal/domain/errors.go` 中定义错误，返回 `DomainError` 类型

**Q: 生产环境怎么部署？**
A:

1. 使用 Docker Compose，配置环境变量
2. 或编译二进制文件 + systemd
3. 前置 Nginx 配置 HTTPS
4. 修改 JWT secret 为强密钥

## 📚 相关文档

- [GORM 文档](https://gorm.io)
- [Gin 文档](https://gin-gonic.com)
- [Casbin 文档](https://casbin.org)

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！
