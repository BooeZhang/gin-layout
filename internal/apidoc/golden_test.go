package apidoc

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gin-layout/config"
)

var update = flag.Bool("update-golden", false, "update golden swagger.json")

func TestGoldenSwaggerJSON(t *testing.T) {
	reg := NewRegistry()

	// ---- Build a representative set of endpoints mirroring admin.go ----

	// Health check (public, query binding)
	reg.Add(&EndpointRecord{
		Method:  "GET",
		Path:    "/health",
		ReqType: typeOfHealthReq(),
		ResType: typeOfHealthRes(),
		Binding: BindingQuery,
		GroupDoc: DocDefaults{
			Visibility: VisibilityPublic,
		},
		RouteDoc: RouteDoc{
			Summary: "健康检查",
			Tags:    []string{"系统"},
		},
	})

	// Auth: login (public, JSON binding)
	reg.Add(&EndpointRecord{
		Method:  "POST",
		Path:    "/api/auth/login",
		ReqType: typeOfLoginReq(),
		ResType: typeOfLoginRes(),
		Binding: BindingJSON,
		GroupDoc: DocDefaults{
			Visibility:    VisibilityPublic,
			DefaultErrors: []ErrorSpec{{HTTPStatus: 400}, {HTTPStatus: 500}},
			Tags:          []string{"认证"},
		},
		RouteDoc: RouteDoc{
			Summary: "登录",
		},
	})

	// Auth: refresh-token (public)
	reg.Add(&EndpointRecord{
		Method:  "POST",
		Path:    "/api/auth/refresh-token",
		ReqType: typeOfRefreshReq(),
		ResType: typeOfTokenRes(),
		Binding: BindingJSON,
		GroupDoc: DocDefaults{
			Visibility:    VisibilityPublic,
			DefaultErrors: []ErrorSpec{{HTTPStatus: 400}, {HTTPStatus: 500}},
			Tags:          []string{"认证"},
		},
		RouteDoc: RouteDoc{
			Summary: "刷新令牌",
		},
	})

	// Users: create (protected, JSON binding)
	reg.Add(&EndpointRecord{
		Method:  "POST",
		Path:    "/api/v1/users",
		ReqType: typeOfCreateUserReq(),
		ResType: typeOfUserRes(),
		Binding: BindingJSON,
		GroupDoc: DocDefaults{
			Security:      []SecurityScheme{BearerAuth},
			Visibility:    VisibilityProtected,
			DefaultErrors: []ErrorSpec{{HTTPStatus: 400}, {HTTPStatus: 401}, {HTTPStatus: 403}, {HTTPStatus: 500}},
			Tags:          []string{"用户"},
		},
		RouteDoc: RouteDoc{
			Summary: "创建用户",
		},
	})

	// Users: list (protected, query binding) — flat struct
	reg.Add(&EndpointRecord{
		Method:  "GET",
		Path:    "/api/v1/users",
		ReqType: typeOfListReq(),
		ResType: typeOfUserListRes(),
		Binding: BindingQuery,
		GroupDoc: DocDefaults{
			Security:      []SecurityScheme{BearerAuth},
			Visibility:    VisibilityProtected,
			DefaultErrors: []ErrorSpec{{HTTPStatus: 400}, {HTTPStatus: 401}, {HTTPStatus: 403}, {HTTPStatus: 500}},
			Tags:          []string{"用户"},
		},
		RouteDoc: RouteDoc{
			Summary: "用户列表",
		},
	})

	// Roles: list (protected, query binding) — generic PageResult
	reg.Add(&EndpointRecord{
		Method:  "GET",
		Path:    "/api/v1/roles",
		ReqType: typeOfListReq(),
		ResType: typeOfGenericPageResult(),
		Binding: BindingQuery,
		GroupDoc: DocDefaults{
			Security:      []SecurityScheme{BearerAuth},
			Visibility:    VisibilityProtected,
			DefaultErrors: []ErrorSpec{{HTTPStatus: 400}, {HTTPStatus: 401}, {HTTPStatus: 403}, {HTTPStatus: 500}},
			Tags:          []string{"角色"},
		},
		RouteDoc: RouteDoc{
			Summary: "角色列表",
		},
	})

	// Users: update (protected, path param, JSON binding)
	recUpdate := &EndpointRecord{
		Method:  "PUT",
		Path:    "/api/v1/users/:id",
		ReqType: typeOfUpdateUserReq(),
		ResType: typeOfUserRes(),
		Binding: BindingJSON,
		GroupDoc: DocDefaults{
			Security:      []SecurityScheme{BearerAuth},
			Visibility:    VisibilityProtected,
			DefaultErrors: []ErrorSpec{{HTTPStatus: 400}, {HTTPStatus: 401}, {HTTPStatus: 403}, {HTTPStatus: 500}},
			Tags:          []string{"用户"},
		},
		RouteDoc: RouteDoc{
			Summary: "更新用户",
			Errors:  []ErrorSpec{{HTTPStatus: 404, Message: "用户不存在"}},
		},
	}
	reg.Add(recUpdate)

	// Users: delete (protected, path param, query binding)
	reg.Add(&EndpointRecord{
		Method:  "DELETE",
		Path:    "/api/v1/users/:id",
		ReqType: typeOfDeleteReq(),
		ResType: typeOfEmptyRes(),
		Binding: BindingQuery,
		GroupDoc: DocDefaults{
			Security:      []SecurityScheme{BearerAuth},
			Visibility:    VisibilityProtected,
			DefaultErrors: []ErrorSpec{{HTTPStatus: 400}, {HTTPStatus: 401}, {HTTPStatus: 403}, {HTTPStatus: 500}},
			Tags:          []string{"用户"},
		},
		RouteDoc: RouteDoc{
			Summary: "删除用户",
		},
	})

	// ---- Build and render ----
	cfg := config.APIDocConfig{
		Enabled:     true,
		Title:       "Gin Layout API",
		Version:     "1.0.0",
		Description: "Gin Layout Swagger API 文档",
	}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	renderer := &Swagger2Renderer{}
	got, err := renderer.Render(spec)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}

	// ---- Compare with golden file ----
	goldenPath := filepath.Join("testdata", "swagger.json")

	if *update {
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatalf("write golden file failed: %v", err)
		}
		t.Log("golden file updated")
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden file failed: %v\nRun 'go test -run TestGoldenSwaggerJSON -update-golden' to generate it.", err)
	}

	if string(got) != string(want) {
		t.Errorf("swagger JSON differs from golden file.\nGot:\n%s\n\nRun 'go test -run TestGoldenSwaggerJSON -update-golden' to update.", string(got))
	}
}

// ---- Helper type constructors (matching real DTO shapes) ----

type goldenHealthRes struct {
	Status string `json:"status"`
}

type goldenLoginReq struct {
	Account  string `json:"account" binding:"required" desc:"账号"`
	Password string `json:"password" binding:"required" desc:"密码"`
}

type goldenLoginRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type goldenRefreshReq struct {
	RefreshToken string `json:"refreshToken" binding:"required" desc:"刷新令牌"`
}

type goldenTokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

type goldenCreateUserReq struct {
	Account  string `json:"account" binding:"required" desc:"账号"`
	Password string `json:"password" binding:"required" desc:"密码"`
	Name     string `json:"name" desc:"姓名"`
}

type goldenUpdateUserReq struct {
	Name string `json:"name" desc:"姓名"`
}

type goldenUserRes struct {
	ID      int64  `json:"id"`
	Account string `json:"account"`
	Name    string `json:"name"`
}

type goldenUserListRes struct {
	Items []goldenUserRes `json:"items"`
	Total int64           `json:"total"`
}

type goldenListReq struct {
	Page int `json:"page" form:"page" desc:"页码"`
	Size int `json:"size" form:"size" desc:"每页数量"`
}

type goldenDeleteReq struct {
	ID uint `json:"id" form:"id" uri:"id"`
}

// goldenPageResult mirrors domain.PageResult for testing generic type doc generation.
type goldenPageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

type goldenRoleItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type goldenHealthReq struct{}

type goldenEmptyRes struct{}

func typeOfHealthReq() reflect.Type     { return reflect.TypeOf(goldenHealthReq{}) }
func typeOfHealthRes() reflect.Type     { return reflect.TypeOf(goldenHealthRes{}) }
func typeOfLoginReq() reflect.Type      { return reflect.TypeOf(goldenLoginReq{}) }
func typeOfLoginRes() reflect.Type      { return reflect.TypeOf(goldenLoginRes{}) }
func typeOfRefreshReq() reflect.Type    { return reflect.TypeOf(goldenRefreshReq{}) }
func typeOfTokenRes() reflect.Type      { return reflect.TypeOf(goldenTokenRes{}) }
func typeOfCreateUserReq() reflect.Type { return reflect.TypeOf(goldenCreateUserReq{}) }
func typeOfUpdateUserReq() reflect.Type { return reflect.TypeOf(goldenUpdateUserReq{}) }
func typeOfUserRes() reflect.Type       { return reflect.TypeOf(goldenUserRes{}) }
func typeOfUserListRes() reflect.Type   { return reflect.TypeOf(goldenUserListRes{}) }
func typeOfGenericPageResult() reflect.Type {
	return reflect.TypeOf(goldenPageResult[goldenRoleItem]{})
}
func typeOfListReq() reflect.Type   { return reflect.TypeOf(goldenListReq{}) }
func typeOfDeleteReq() reflect.Type { return reflect.TypeOf(goldenDeleteReq{}) }
func typeOfEmptyRes() reflect.Type  { return reflect.TypeOf(goldenEmptyRes{}) }
