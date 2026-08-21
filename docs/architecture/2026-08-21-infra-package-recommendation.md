# internal/infra 包架构建议：合并、拆分还是删除？

> 日期：2026-08-21
> 范围：`internal/infra` 及其与 `internal/token`、`internal/policy`、`internal/page`、`internal/reqctx` 的边界
> 结论：**当前规模下保持现状；触发信号出现后按技术拆分至 `infra` 消失；合并与"删除后塞回领域包"均不可取。**

---

## 1. 背景与问题

针对 `internal/infra` 的定位，本仓库先后提出了四个连续问题：

1. `internal/token` 中的 `Issuer` 接口是否多余？
2. `internal/token` 是否应移入 `internal/infra`？
3. `internal/page`、`internal/policy`、`internal/reqctx` 是否都可移入 `internal/infra`？
4. `internal/infra` 究竟该合并、拆分还是删除？

本文汇总分析结论、关键依据、目标架构与迁移方案，供后续决策使用。

---

## 2. 现状盘点

### 2.1 `internal/infra` 现有内容

| 文件 | 内容 | 类别 |
| --- | --- | --- |
| `cache.go` | `RedisClient`：单机/集群模式选择、`Ping`/`Close` | 通用基建 |
| `database.go` | `Database`：多驱动 GORM 连接、连接池、`Migrate`/`Ping`/`Close`、`ErrDbIsNil` | 通用基建 |
| `crud.go` | `CRUDRepository[M, ID]`：泛型 CRUD 基座（`Create`/`Update`/`Delete`/`FindByID`/`FindByIDs`） | 通用基建 |
| `logger.go` | `Logger = zerolog.Logger` 别名、`NewLogger`、`DefaultLogger`、`LogFromContext` | 通用基建 |
| `jwt.go` | `JWTIssuer`：JWT 签发/解析、错误归一化，**实现 `token.Issuer`** | 契约实现 |
| `policy.go` | `CasbinManager`：角色/权限同步、`Enforce`，**实现 `policy.Manager`** | 契约实现 |
| `token_blacklist.go` | `TokenBlacklistModel` + `TokenBlacklistRepository`，**实现 `token.BlacklistRepository`** | 契约实现 |

### 2.2 依赖关系（`go list` 实测）

```
page    -> []                                          （纯叶子）
reqctx  -> [context]                                   （纯叶子）
policy  -> [context, gin-layout/internal/apperror]     （近叶子）
infra   -> [config, reqctx, token, casbin, gorm, redis,
            zerolog, jwt, lo, ...]                     （依赖第三方与技术栈）
```

关键点：
- `infra` 依赖 `token`（`jwt.go` 使用 `token.Pair`/`Claims`/错误），但**不依赖 `policy`**——`CasbinManager` 通过 Go 隐式接口满足 `policy.Manager`。
- `infra` 依赖 `reqctx`（`logger.go` 从 context 提取请求级字段）。

### 2.3 `infra` 的消费者

| 消费者 | 使用的 `infra` 符号 |
| --- | --- |
| `bootstrap/app.go`、`bootstrap/initializer` | 组装：`NewDatabase`/`NewRedis`/`NewJWTIssuer`/`NewCasbinManager`/`TokenBlacklistRepository`、`Logger` |
| `middleware/recovery.go`、`logger.go`、`error.go` | `Logger`、`LogFromContext` |
| `server/server.go` | `Logger` |
| `migration/migration.go`、`cmd/migrate/main.go` | `Database`、`ErrDbIsNil`、`TokenBlacklistModel`、`NewLogger` |
| `sysuser/repo.go`、`role/repo.go`、`menu/repo.go` | `CRUDRepository` |
| `common/base.go` | `LogFromContext` |

---

## 3. 结论摘要

| 问题 | 结论 | 一句话理由 |
| --- | --- | --- |
| `Issuer` 是否多余 | **否** | 是打破 `infra → token` 循环的必需接缝，且定义在使用方 |
| `token` 移入 `infra` | **否** | 会把业务逻辑埋进基础设施，业务包被迫依赖实现细节 |
| `page`/`policy`/`reqctx` 移入 `infra` | **否** | `reqctx` 有硬性循环约束；`policy` 合并会破坏"接口在消费方" |
| `infra` 合并 | **否** | 合并方向与前三项相同，只会让大杂烩更大 |
| `infra` 拆分 | **触发信号后** | 收益当前是"美观税"，规模变大后再做 |
| `infra` 删除 | **仅作为拆分完成后的终点** | 直接删除且不拆分 = 把实现塞回契约包，破坏隐式接口 |

---

## 4. 关键论证

### 4.1 分层心智模型

本仓库的包应分为四类：

1. **叶子共享词汇**（无内部依赖）：`apperror`（错误模型）、`page`（分页类型）、`reqctx`（请求上下文）。
2. **消费方契约**（定义接口、不背持久化）：`token`（`Issuer`/`BlacklistRepository`/`Manager`）、`policy`（`Manager`/`PermissionResolver`）。
3. **技术实现**（有状态、有外部依赖）：`infra` 现有全部内容。
4. **组装**：`bootstrap` 把实现注入契约。

`infra` 的问题不是"存在"，而是**同时容纳了第 3 类中的两种东西**：通用基建（db/redis/crud/logger）与契约实现（jwt/casbin/blacklist），并且按"层"命名而非按"职责"命名。

### 4.2 为什么 `reqctx` 不能进 `infra`（循环证明）

```
token/token.go     依赖 reqctx（CurrentUserFromContext / CurrentTokenFromContext）
infra/jwt.go       依赖 token（Pair / Claims / 错误）
```

若 `reqctx` 并入 `infra`：`token → infra → token`，循环导入，编译失败。`reqctx` 必须是叶子包。

### 4.3 为什么 `policy` 不能进 `infra`

- `policy.Manager` 定义在消费方一侧，`infra.CasbinManager` 隐式实现且不 import `policy`——这正是 Go"接口定义在使用方"的惯用法。
- `policy` 的消费者（`middleware/rbac.go`、`sysuser/service.go`、`role/service.go`、`router/admin.go`、`web/errors_test.go`）共 7 处。若并入 `infra`，这些业务包将被强制 import 携带 casbin/gorm/redis 的 `infra`，实现细节扩散到业务面。

### 4.4 为什么"删除并塞回领域包"不行

直接删除 `infra` 而不拆分，`JWTIssuer`/`CasbinManager` 无处安放；塞回 `token`/`policy` 会让契约包背上 jwt/casbin/gorm 依赖，破坏隐式接口接缝。删除的**唯一正确路径**是先拆分、后删除。

### 4.5 为什么现在保持现状

- 拆分收益当前是"美观税"：7 个文件各自内聚、每个契约只有一个实现、单一二进制、消费方约 10 个文件 20 余行改动，行为零变化。
- `gin-layout` 是模板项目，README 已写明"分层架构 + `infra/` = 基础设施层"，`infra` 目录本身是模板约定的一部分，团队认知成本低。
- 唯一的低成本动作（可选）：把 `jwt.go` + `jwt_test.go` 独立成 `internal/jwt`，改动仅 `bootstrap/app.go` 一处。

---

## 5. 目标架构（触发拆分后）

```
internal/
├── apperror/        # 错误模型（叶子）
├── page/            # 分页类型（叶子）
├── reqctx/          # 请求上下文（叶子）
├── token/           # 契约：Issuer / BlacklistRepository / Manager + 编排（可选自持黑名单持久化）
├── policy/          # 契约：Manager / PermissionResolver
├── jwt/             # 实现 token.Issuer（保持 jwt → token 方向）
├── casbin/          # 实现 policy.Manager
├── db/              # GORM 连接 + 泛型 CRUD 基座（database.go + crud.go + ErrDbIsNil）
├── redis/           # Redis 连接（cache.go）
├── logger/          # 日志封装（logger.go）
├── sysuser/ role/ menu/ health/   # 业务领域（含各自 repo）
├── middleware/ router/ web/ server/ bootstrap/ config/ migration/
```

文件映射：

| `infra` 现有文件 | 目标包 | 消费者改动 |
| --- | --- | --- |
| `cache.go` | `internal/redis` | `bootstrap` |
| `database.go` + `crud.go` | `internal/db` | `migration`、`cmd/migrate`、`sysuser`/`role`/`menu` repo、`bootstrap` |
| `logger.go` | `internal/logger` | `middleware`×3、`server`、`common/base`、`bootstrap`、`initializer` |
| `jwt.go` + `jwt_test.go` | `internal/jwt` | `bootstrap` |
| `policy.go` | `internal/casbin` | `bootstrap` |
| `token_blacklist.go` | `internal/token`（推荐）或独立 `internal/tokenblacklist` | `migration`、`bootstrap` |

黑名单归属的取舍：
- **方案 A（推荐）**：`TokenBlacklistRepository` 移入 `internal/token`，`Add` 改为直接使用 `*gorm.DB`（不再依赖 `db.CRUDRepository`）。与本仓库 `sysuser`/`role`/`menu` 领域包自持持久化的约定一致，且无循环（`token` 只依赖 gorm 第三方库）。
- **方案 B（保守）**：保留独立 `internal/tokenblacklist` 包，继续复用 `db.CRUDRepository`。适合不希望 `token` 包出现 gorm 依赖的情况。

---

## 6. 决策信号（何时动手拆分）

满足任一信号即可启动第 7 节迁移：

- 出现**第二个实现**：如 Redis 版黑名单、Paseto 版签发——接口的价值只有在多实现时体现。
- 引入**新第三方基建**（Kafka、S3、消息队列等），需要新包而非继续堆进 `infra`。
- `infra` 文件数明显增长（如超过 15 个）或出现命名冲突。
- 团队在 review 中频繁出现"不要往 `infra` 里加"的讨论，或成员经常找不到文件。

---

## 7. 迁移方案（若执行拆分）

### 阶段 1：拆分通用基建（无契约依赖，风险最低）

1. 创建 `internal/db`：迁入 `database.go` + `crud.go`（含 `ErrDbIsNil`）。
   - 改动：`migration/migration.go`、`cmd/migrate/main.go`、`sysuser/repo.go`、`role/repo.go`、`menu/repo.go`、`bootstrap/app.go`。
2. 创建 `internal/redis`：迁入 `cache.go`。
   - 改动：`bootstrap/app.go`。
3. 创建 `internal/logger`：迁入 `logger.go`。
   - 改动：`middleware/recovery.go`、`middleware/logger.go`、`middleware/error.go`、`server/server.go`、`common/base.go`、`bootstrap/app.go`、`bootstrap/initializer/initializer.go`。

### 阶段 2：拆分契约实现

4. 创建 `internal/jwt`：迁入 `jwt.go` + `jwt_test.go`，保持 `jwt → token` 依赖方向。
   - 改动：`bootstrap/app.go`（`infra.NewJWTIssuer` → `jwt.NewJWTIssuer`）。
5. 创建 `internal/casbin`：迁入 `policy.go`，保持不依赖 `policy` 包的隐式实现。
   - 改动：`bootstrap/app.go`。

### 阶段 3：黑名单归属

6. 按方案 A：`token_blacklist.go` 移入 `internal/token`，`Add` 直接使用 `*gorm.DB`。
   - 改动：`migration/migration.go`（引用 `token.TokenBlacklistModel`）、`bootstrap/app.go`。
   - 或按方案 B：保留独立 `internal/tokenblacklist` 包。

### 阶段 4：删除 `infra` 并验证

7. 删除 `internal/infra`，全仓替换残留引用。
8. 验证：`go build ./...`、`go vet ./...`、`go test ./...`。

改动量估算：约 10–12 个文件、20–30 行 import/引用变化，无行为变化。

---

## 8. 一句话结论

**当前规模下，合并、拆分、删除都无正收益——保持现状；以第 5 节为目标、第 6 节为触发信号，规模扩大后按第 7 节一次性拆分，届时 `infra` 包自然消失。**
