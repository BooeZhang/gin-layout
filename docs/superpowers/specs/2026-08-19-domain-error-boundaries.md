# Domain Error Boundary Refactor Spec

Date: 2026-08-19
Project: gin-layout

## Problem

`internal/domain/error.go` currently makes the domain package aware of HTTP status codes, response messages, and a response-oriented error interface. `internal/domain/errors.go` also keeps generic not-found / invalid-argument sentinels at the wrong layer. On top of that, `internal/infra/errors.go` translates storage-layer not-found into a domain-layer not-found sentinel, so persistence semantics leak upward.

That shape is not aligned with DDD + clean architecture:

- the domain layer should not know about transport concerns;
- repository code should not depend on domain-level response contracts;
- the web adapter should own HTTP status, business code, and response message translation.

## Goals

- Keep `internal/domain` focused on entities, value objects, and pure domain behavior.
- Move every business error sentinel to the package that owns the behavior.
- Make single-record repository lookups report miss/hit explicitly instead of using a shared not-found error.
- Keep `internal/web` as the only place that maps internal errors to HTTP status, business code, and message.
- Preserve the current response envelope shape: `{code, message, data}`.
- Keep `errors.Is` as the only matching mechanism for sentinel errors.

## Non-Goals

- Do not move `SysUser`, `Role`, or `Menu` out of `internal/domain`.
- Do not introduce a new framework, error library, or layered `common/errors` package.
- Do not redesign auth, routing, or frontend behavior as part of this refactor.
- Do not change the JSON envelope shape.

## Target State

### Error ownership

The following packages own their own sentinels and import only the standard library for error creation:

- `internal/sysuser`
- `internal/role`
- `internal/menu`
- `internal/token`
- `internal/policy`

### Repository contract

Every repository method that returns one logical row must return an explicit existence flag:

```go
(*T, bool, error)
```

Rules:

- `found == false` means the row did not exist.
- `found == false` is not an error.
- storage-layer not-found must never be converted into a domain sentinel.

### Web translation

`internal/web` owns a local `ErrorDescriptor` and a `DecodeError` function that converts known sentinels into:

- HTTP status
- response `code`
- response `message`

Unknown errors map to a single internal-server-error descriptor.

## Canonical Error Vocabulary

Preserve the current outward code/message values unless a task explicitly says otherwise. The following vocabulary is the intended final set:

| Package | Sentinel | HTTP | Code | Message |
| --- | --- | ---: | ---: | --- |
| `sysuser` | `ErrInvalidAccountFormat` | 422 | 20010 | `用户名格式错` |
| `sysuser` | `ErrInvalidUserID` | 422 | 20011 | `无效的用户 ID` |
| `sysuser` | `ErrWeakPassword` | 422 | 20012 | `密码强度不足` |
| `sysuser` | `ErrAccountExists` | 409 | 20020 | `账号已存在` |
| `sysuser` | `ErrUserNotFound` | 404 | 20000 | `用户不存在` |
| `sysuser` | `ErrUserDisabled` | 200 | 20040 | `用户已禁用` |
| `sysuser` | `ErrCannotDeleteAdmin` | 200 | 20050 | `不能删除超级管理员` |
| `sysuser` | `ErrInvalidCredentials` | 200 | 20051 | `用户名或密码错误` |
| `sysuser` | `ErrPasswordIdentical` | 200 | 20052 | `新密码不能与旧密码相同` |
| `role` | `ErrInvalidRoleID` | 422 | 30010 | `无效ID` |
| `role` | `ErrRoleExists` | 409 | 30020 | `角色已存在` |
| `role` | `ErrRoleNotFound` | 404 | 30000 | `角色不存在` |
| `role` | `ErrPermissionNotFound` | 404 | 30100 | `权限不存在` |
| `role` | `ErrRoleDisabled` | 200 | 30040 | `角色已禁用` |
| `role` | `ErrCannotDeleteAdminRole` | 200 | 30050 | `不允许删除管理员角色` |
| `menu` | `ErrInvalidMenuID` | 422 | 40010 | `无效ID` |
| `menu` | `ErrMenuExists` | 409 | 40020 | `菜单已存在` |
| `menu` | `ErrMenuNotFound` | 404 | 40000 | `菜单不存在` |
| `token` | `ErrInvalidAccessToken` | 401 | 50010 | `无效访问令牌` |
| `token` | `ErrUnauthenticated` | 401 | 50011 | `未登录或非法访问` |
| `token` | `ErrTokenInvalid` | 401 | 50012 | `token 无效` |
| `token` | `ErrTokenExpired` | 401 | 50060 | `token 已过期` |
| `token` | `ErrTokenRevoked` | 401 | 50061 | `token 已失效` |
| `token` | `ErrTokenNotActive` | 401 | 50070 | `token 不是活跃状态` |
| `policy` | `ErrPermissionDenied` | 403 | 30130 | `没有权限` |

## Validation Criteria

- `rg -n "DomainError|CodedError|NewDomainError|domain\\.ErrNotFound|domain\\.ErrInvalidArgument|NormalizeError" internal` returns no matches.
- `go test ./...` passes.
- Wrapped sentinel errors still match with `errors.Is`.

