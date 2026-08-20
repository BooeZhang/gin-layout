# Domain Error Boundary Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove transport and persistence error semantics from `internal/domain`, move error ownership to the packages that own the behavior, and keep HTTP/status/message translation inside `internal/web`.

**Architecture:** Keep `internal/domain` entity-only for this refactor. Each feature package (`sysuser`, `role`, `menu`, `token`, `policy`) owns its own sentinel errors with `errors.New` and matches them with `errors.Is`. Repository methods report hit/miss explicitly for single-row lookups instead of converting storage misses into domain errors. `internal/web` becomes the single translation point from internal errors to HTTP status, response code, and message.

**Tech Stack:** Go 1.26.3, stdlib `errors`, GORM, gin, table-driven tests, in-memory sqlite for repository verification.

**Spec:** `docs/superpowers/specs/2026-08-19-domain-error-boundaries.md`

## Global Constraints

- Keep the response envelope shape `{code, message, data}` unchanged.
- Do not reintroduce `DomainError`, `CodedError`, `NewDomainError`, `domain.ErrNotFound`, or `domain.ErrInvalidArgument`.
- Keep `internal/domain` free of `net/http`.
- Match errors only with `errors.Is` against package-local sentinels or wrapped sentinels.
- Single-record repository lookups must return `(*T, bool, error)`.
- Preserve the current outward code/message values listed in the spec.

---

### Task 1: Introduce package-local error vocabularies and web decoding

**Files:**
- Modify: `internal/sysuser/errors.go`
- Modify: `internal/role/errors.go`
- Modify: `internal/menu/errors.go`
- Modify: `internal/token/token.go`
- Modify: `internal/policy/policy.go`
- Create: `internal/web/errors.go`
- Modify: `internal/web/response.go`
- Modify: `internal/middleware/error.go`

**Interfaces:**
- Consumes: `errors.New`, `errors.Is`, current code/message literals.
- Produces: `sysuser.ErrUserNotFound`, `sysuser.ErrInvalidCredentials`, `role.ErrRoleNotFound`, `role.ErrPermissionNotFound`, `menu.ErrMenuNotFound`, `token.ErrUnauthenticated`, `policy.ErrPermissionDenied`, `web.ErrorDescriptor`, `web.DecodeError(err error) ErrorDescriptor`.

- [ ] **Step 1: Add the tests that define the new web translation contract**

Create `internal/web/errors_test.go` with a table that asserts direct and wrapped sentinels decode to the current outward tuple.

```go
cases := []struct {
	name string
	err  error
	want ErrorDescriptor
}{
	{
		name: "sysuser not found",
		err:  sysuser.ErrUserNotFound,
		want: ErrorDescriptor{HTTPStatus: 404, Code: 20000, Message: "用户不存在"},
	},
	{
		name: "wrapped unauthenticated",
		err:  fmt.Errorf("wrap: %w", token.ErrUnauthenticated),
		want: ErrorDescriptor{HTTPStatus: 401, Code: 50011, Message: "未登录或非法访问"},
	},
}
```

- [ ] **Step 2: Replace the current domain-backed sentinels with package-local sentinels**

In each owning package, define the final sentinel variables with `errors.New(...)` and keep the current outward messages attached to the same behavior.

- [ ] **Step 3: Add `internal/web/errors.go`**

Define:

```go
type ErrorDescriptor struct {
	HTTPStatus int
	Code       int
	Message    string
}

func DecodeError(err error) ErrorDescriptor
```

`DecodeError` must match by `errors.Is`, not by concrete type assertions.

- [ ] **Step 4: Update `internal/web/response.go` and `internal/middleware/error.go` to use the new decoder**

`web.Error` should render the decoded descriptor. `middleware.ErrorHandler` should log the decoded descriptor fields instead of asserting `domain.DomainError`.

- [ ] **Step 5: Run the focused package tests**

Run:

```bash
go test ./internal/web ./internal/middleware ./internal/token ./internal/policy ./internal/sysuser ./internal/role ./internal/menu
```

Expected: pass after the call sites are migrated in later tasks, or compile-only failures only if Task 2 has not landed yet.

### Task 2: Change single-record repository methods to report hit/miss explicitly

**Files:**
- Modify: `internal/infra/crud.go`
- Modify: `internal/infra/token_blacklist.go`
- Modify: `internal/sysuser/ports.go`
- Modify: `internal/role/ports.go`
- Modify: `internal/menu/ports.go`
- Modify: `internal/sysuser/repo.go`
- Modify: `internal/role/repo.go`
- Modify: `internal/menu/repo.go`

**Interfaces:**
- Consumes: `gorm.ErrRecordNotFound`, `gorm.G[...]().First(...)`, current repository ports.
- Produces: `FindByID(ctx, id) (*T, bool, error)`, `FindByAccount(ctx, account) (*domain.SysUser, bool, error)`, `FindByCode(ctx, code) (*domain.Role, bool, error)`, `FindByCode(ctx, code) (*domain.Menu, bool, error)`, `FindByIDWithRoles(ctx, id) (*domain.SysUser, bool, error)`, `FindByIDWithPerm(ctx, roleID) (*domain.Role, bool, error)`.

- [ ] **Step 1: Write a failing sqlite-backed repository test**

Create `internal/infra/crud_test.go` and verify that a missing row returns `found=false` and `err=nil`.

```go
got, found, err := repo.FindByID(ctx, int64(999))
if err != nil {
	t.Fatalf("unexpected error: %v", err)
}
if found {
	t.Fatalf("expected not found")
}
if got != nil {
	t.Fatalf("expected nil entity")
}
```

- [ ] **Step 2: Change `CRUDRepository.FindByID`**

Return the entity plus a `found` boolean. When GORM reports record-not-found, return `nil, false, nil`. For all other errors, return `nil, false, err`.

- [ ] **Step 3: Update the concrete repositories**

Apply the same contract to the feature repositories that currently expose single-row lookup methods:

```go
FindByAccount
FindByCode
FindByIDWithRoles
FindByIDWithPerm
```

- [ ] **Step 4: Remove `internal/infra/errors.go` usage**

Delete the `NormalizeError` call sites from the repo implementations and from `internal/infra/token_blacklist.go`. Missing rows are no longer normalized into a domain error.

- [ ] **Step 5: Run the focused repository tests**

Run:

```bash
go test ./internal/infra ./internal/sysuser ./internal/role ./internal/menu
```

Expected: the repository packages compile against the new contract, even if service call sites still need Task 3.

### Task 3: Migrate services, bootstrap, and auth/policy call sites to the new contract

**Files:**
- Modify: `internal/sysuser/service.go`
- Modify: `internal/role/service.go`
- Modify: `internal/menu/service.go`
- Modify: `internal/bootstrap/initializer/initializer.go`
- Modify: `internal/token/token.go`
- Modify: `internal/policy/policy.go`
- Modify: `internal/middleware/error.go`

**Interfaces:**
- Consumes: `(*T, bool, error)` repository methods, package-local not-found sentinels, `web.DecodeError`.
- Produces: service-level branches that return `ErrUserNotFound`, `ErrRoleNotFound`, `ErrMenuNotFound`, `ErrInvalidCredentials`, `ErrUnauthenticated`, and `ErrPermissionDenied` without any `domain.ErrNotFound` dependency.

- [ ] **Step 1: Replace every `errors.Is(err, domain.ErrNotFound)` branch**

Use the `found` boolean from repositories where the lookup is performed inside the current package. Only keep `errors.Is` for wrapped local sentinels that are already part of the package vocabulary.

Concrete pattern:

```go
entity, found, err := s.repo.FindByID(ctx, id)
if err != nil {
	return nil, err
}
if !found {
	return nil, ErrRoleNotFound
}
```

- [ ] **Step 2: Rename the remaining exported error vars to the final vocabulary**

Use the names from the spec exactly: `ErrUserNotFound`, `ErrRoleNotFound`, `ErrMenuNotFound`, `ErrInvalidCredentials`, `ErrUnauthenticated`, and `ErrPermissionDenied`.

- [ ] **Step 3: Update `internal/bootstrap/initializer/initializer.go`**

The bootstrapper must treat missing roles and users as package-local not-found cases and create the seed data only when `found == false`.

- [ ] **Step 4: Update auth-related flows**

`internal/token/token.go` must return `ErrUnauthenticated` for the “no current user / no current token” cases. `internal/policy/policy.go` must expose `ErrPermissionDenied` as the package-local sentinel.

- [ ] **Step 5: Run the targeted service tests**

Run:

```bash
go test ./internal/bootstrap/initializer ./internal/sysuser ./internal/role ./internal/menu ./internal/token ./internal/policy
```

Expected: all call sites now compile against the new repository contract and package-local sentinels.

### Task 4: Remove the legacy domain and infra error files

**Files:**
- Delete: `internal/domain/error.go`
- Delete: `internal/domain/errors.go`
- Delete: `internal/infra/errors.go`

**Interfaces:**
- Consumes: no remaining references to `domain.DomainError`, `domain.CodedError`, `domain.ErrNotFound`, `domain.ErrInvalidArgument`, or `infra.NormalizeError`.
- Produces: a domain package that no longer knows about HTTP status, response codes, or shared not-found sentinels.

- [ ] **Step 1: Search for stale references before deleting**

Run:

```bash
rg -n "DomainError|CodedError|NewDomainError|domain\\.ErrNotFound|domain\\.ErrInvalidArgument|NormalizeError" internal
```

Any remaining matches must be resolved before the files are deleted.

- [ ] **Step 2: Delete the three legacy files**

Remove the old domain error model and the infra normalization shim.

- [ ] **Step 3: Fix any compiler fallout**

Remove stale imports and any leftover code paths that still expect the old error model.

- [ ] **Step 4: Run the full Go test suite**

Run:

```bash
go test ./...
```

Expected: the repository builds and tests with no dependency on the deleted error layer.

### Task 5: Harden verification with focused tests

**Files:**
- Create or modify: `internal/web/errors_test.go`
- Create or modify: `internal/infra/crud_test.go`
- Create or modify: `internal/sysuser/service_test.go`
- Create or modify: `internal/role/service_test.go`
- Create or modify: `internal/menu/service_test.go`
- Create or modify: `internal/token/token_test.go`
- Create or modify: `internal/policy/policy_test.go`

**Interfaces:**
- Consumes: the final sentinel vocabulary, the `(*T, bool, error)` repository contract, and `web.DecodeError`.
- Produces: table-driven tests that prove the new boundary stays stable.

- [ ] **Step 1: Verify wrapped sentinel matching**

Every package-local sentinel that is intended to be public must match through `fmt.Errorf("wrap: %w", err)` via `errors.Is`.

- [ ] **Step 2: Verify repository hit/miss behavior**

Add tests for:

```go
found == true, err == nil
found == false, err == nil
err != nil
```

for the single-row lookup methods.

- [ ] **Step 3: Verify service translation**

Add tests that prove each service converts `found == false` into its own not-found sentinel and does not leak repository semantics upward.

- [ ] **Step 4: Run the final verification commands**

Run:

```bash
go test ./...
rg -n "DomainError|CodedError|NewDomainError|domain\\.ErrNotFound|domain\\.ErrInvalidArgument|NormalizeError" internal
```

Expected:

- `go test ./...` passes.
- The `rg` command returns no matches.

