package apidoc

import (
	"reflect"
	"strings"
	"testing"

	"gin-layout/config"
)

func TestDefaultSummary(t *testing.T) {
	tests := []struct {
		method, path, want string
	}{
		{"POST", "/users", "Create users"},
		{"GET", "/users", "List users"},
		{"GET", "/users/:id", "Get users"},
		{"PUT", "/users/:id", "Update users"},
		{"DELETE", "/users/:id", "Delete users"},
		{"GET", "/health", "List health"},
	}

	for _, tt := range tests {
		got := defaultSummary(tt.method, tt.path)
		if got != tt.want {
			t.Errorf("defaultSummary(%q, %q) = %q, want %q", tt.method, tt.path, got, tt.want)
		}
	}
}

func TestParsePathParams(t *testing.T) {
	params := parsePathParams("/api/v1/users/:userId/roles/:roleId")
	if len(params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(params))
	}

	if params[0].Name != "userId" || params[0].In != "path" {
		t.Errorf("first param: want userId/path, got %s/%s", params[0].Name, params[0].In)
	}
	if params[1].Name != "roleId" || params[1].In != "path" {
		t.Errorf("second param: want roleId/path, got %s/%s", params[1].Name, params[1].In)
	}
	if !params[0].Required {
		t.Error("path params should be required")
	}
	if params[0].Type != "string" {
		t.Errorf("path params should default to string, got %q", params[0].Type)
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"/api/users/:id", "/api/users/{id}"},
		{"/api/v1/users/:userId/roles/:roleId", "/api/v1/users/{userId}/roles/{roleId}"},
		{"/health", "/health"},
		{"api/users", "/api/users"},
	}

	for _, tt := range tests {
		got := normalizePath(tt.in)
		if got != tt.want {
			t.Errorf("normalizePath(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestBuilder_PublicRoute(t *testing.T) {
	reg := NewRegistry()
	type EmptyReq struct{}
	type HealthRes struct {
		Status string `json:"status"`
	}

	// 模拟路由注册：没有安全认证的公开路由。
	rec := &EndpointRecord{
		Method:  "GET",
		Path:    "/health",
		ReqType: reflect.TypeOf(EmptyReq{}),
		ResType: reflect.TypeOf(HealthRes{}),
		Binding: BindingQuery,
		GroupDoc: DocDefaults{
			Visibility: VisibilityPublic,
		},
	}
	reg.Add(rec)

	cfg := config.APIDocConfig{
		Enabled: true,
		Title:   "Test API",
		Version: "1.0",
	}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	if spec.Title != "Test API" {
		t.Errorf("expected title 'Test API', got %q", spec.Title)
	}

	// 验证路径存在。
	swPath := normalizePath("/health")
	pathModel, ok := spec.Paths[swPath]
	if !ok {
		t.Fatalf("expected path %q in spec", swPath)
	}

	op, ok := pathModel["get"]
	if !ok {
		t.Fatal("expected GET operation")
	}

	// 公开路由不应有安全配置。
	if len(op.Security) != 0 {
		t.Errorf("public route should have no security, got %v", op.Security)
	}

	// 应有 200 响应。
	if _, ok := op.Responses["200"]; !ok {
		t.Error("expected 200 response")
	}

	// 应有默认摘要。
	if op.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestBuilder_ProtectedRoute(t *testing.T) {
	reg := NewRegistry()
	type CreateUserReq struct {
		Name string `json:"name" binding:"required"`
	}
	type UserRes struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	rec := &EndpointRecord{
		Method:  "POST",
		Path:    "/api/v1/users",
		ReqType: reflect.TypeOf(CreateUserReq{}),
		ResType: reflect.TypeOf(UserRes{}),
		Binding: BindingJSON,
		GroupDoc: DocDefaults{
			Security:   []SecurityScheme{BearerAuth},
			Visibility: VisibilityProtected,
			Tags:       []string{"用户"},
		},
	}
	reg.Add(rec)

	cfg := config.APIDocConfig{Enabled: true, Title: "Test API", Version: "1.0"}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	swPath := normalizePath("/api/v1/users")
	pathModel := spec.Paths[swPath]
	op := pathModel["post"]

	// 受保护路由应有安全配置。
	if len(op.Security) == 0 {
		t.Error("protected route should have security")
	}

	// 应有来自组默认值的标签。
	if len(op.Tags) == 0 {
		t.Error("expected tags")
	}

	// 应有请求体参数。
	hasBody := false
	for _, p := range op.Parameters {
		if p.In == "body" {
			hasBody = true
			if !strings.Contains(p.SchemaRef, "definitions") {
				t.Errorf("expected definitions ref in schema, got %q", p.SchemaRef)
			}
		}
	}
	if !hasBody {
		t.Error("expected body parameter for JSON binding")
	}

	// 响应应有真实的数据 schema（非不透明对象）。
	resp := op.Responses["200"]
	if resp.SchemaRef == "#/definitions/common.Response" {
		t.Error("response schema should include concrete data type, not opaque common.Response")
	}
	if !strings.Contains(resp.SchemaRef, "UserRes") {
		t.Errorf("response schema ref should reference UserRes, got %q", resp.SchemaRef)
	}
}

func TestBuilder_QueryBinding(t *testing.T) {
	reg := NewRegistry()
	type ListReq struct {
		Page int `json:"page" form:"page" desc:"页码"`
		Size int `json:"size" form:"size" desc:"每页数量"`
	}
	type ListRes struct {
		Total int64 `json:"total"`
	}

	rec := &EndpointRecord{
		Method:  "GET",
		Path:    "/api/v1/users",
		ReqType: reflect.TypeOf(ListReq{}),
		ResType: reflect.TypeOf(ListRes{}),
		Binding: BindingQuery,
	}
	reg.Add(rec)

	cfg := config.APIDocConfig{Enabled: true, Title: "Test API", Version: "1.0"}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	swPath := normalizePath("/api/v1/users")
	pathModel := spec.Paths[swPath]
	op := pathModel["get"]

	// 应有查询参数。
	hasPage := false
	hasSize := false
	for _, p := range op.Parameters {
		if p.In == "query" {
			if p.Name == "page" {
				hasPage = true
				if p.Description != "页码" {
					t.Errorf("expected '页码' description, got %q", p.Description)
				}
			}
			if p.Name == "size" {
				hasSize = true
			}
		}
	}
	if !hasPage {
		t.Error("expected 'page' query parameter")
	}
	if !hasSize {
		t.Error("expected 'size' query parameter")
	}
}

func TestBuilder_PathParams(t *testing.T) {
	reg := NewRegistry()
	type EmptyReq struct{}

	rec := &EndpointRecord{
		Method:  "GET",
		Path:    "/api/v1/users/:id",
		ReqType: reflect.TypeOf(EmptyReq{}),
		ResType: reflect.TypeOf(EmptyReq{}),
		Binding: BindingQuery,
	}
	reg.Add(rec)

	cfg := config.APIDocConfig{Enabled: true, Title: "Test API", Version: "1.0"}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	pathKey := "/api/v1/users/{id}"
	if _, ok := spec.Paths[pathKey]; !ok {
		t.Fatalf("expected normalized path %q, got paths: %v", pathKey, spec.Paths)
	}

	op := spec.Paths[pathKey]["get"]
	hasPathParam := false
	for _, p := range op.Parameters {
		if p.In == "path" && p.Name == "id" {
			hasPathParam = true
			if !p.Required {
				t.Error("path param should be required")
			}
		}
	}
	if !hasPathParam {
		t.Error("expected path parameter 'id'")
	}
}

func TestBuilder_ErrorMerging(t *testing.T) {
	reg := NewRegistry()
	type EmptyReq struct{}

	rec := &EndpointRecord{
		Method:  "PUT",
		Path:    "/api/v1/users/:id",
		ReqType: reflect.TypeOf(EmptyReq{}),
		ResType: reflect.TypeOf(EmptyReq{}),
		Binding: BindingJSON,
		GroupDoc: DocDefaults{
			DefaultErrors: []ErrorSpec{
				{HTTPStatus: 400, Message: "请求参数错误"},
				{HTTPStatus: 401, Message: "未授权"},
				{HTTPStatus: 500, Message: "服务器错误"},
			},
		},
		RouteDoc: RouteDoc{
			Errors: []ErrorSpec{
				{HTTPStatus: 404, Message: "用户不存在"},
			},
		},
	}
	reg.Add(rec)

	cfg := config.APIDocConfig{Enabled: true, Title: "Test API", Version: "1.0"}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	op := spec.Paths["/api/v1/users/{id}"]["put"]

	// 应包含 400、401、404、500 状态码。
	for _, status := range []string{"400", "401", "404", "500"} {
		if _, ok := op.Responses[status]; !ok {
			t.Errorf("expected response status %s", status)
		}
	}

	// 404 应有自定义消息。
	if op.Responses["404"].Description != "用户不存在" {
		t.Errorf("expected '用户不存在', got %q", op.Responses["404"].Description)
	}
}

func TestBuilder_HiddenRoute(t *testing.T) {
	reg := NewRegistry()
	type EmptyReq struct{}

	hidden := true
	rec := &EndpointRecord{
		Method:  "GET",
		Path:    "/internal/debug",
		ReqType: reflect.TypeOf(EmptyReq{}),
		ResType: reflect.TypeOf(EmptyReq{}),
		Binding: BindingQuery,
		RouteDoc: RouteDoc{
			Hidden: &hidden,
		},
	}
	reg.Add(rec)

	cfg := config.APIDocConfig{Enabled: true, Title: "Test API", Version: "1.0"}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	if _, ok := spec.Paths["/internal/debug"]; ok {
		t.Error("hidden route should not appear in spec")
	}
}

func TestBuilder_SecurityDefinitions(t *testing.T) {
	reg := NewRegistry()
	type EmptyReq struct{}

	rec := &EndpointRecord{
		Method:  "GET",
		Path:    "/api/v1/protected",
		ReqType: reflect.TypeOf(EmptyReq{}),
		ResType: reflect.TypeOf(EmptyReq{}),
		Binding: BindingQuery,
		GroupDoc: DocDefaults{
			Security:   []SecurityScheme{BearerAuth},
			Visibility: VisibilityProtected,
		},
	}
	reg.Add(rec)

	cfg := config.APIDocConfig{Enabled: true, Title: "Test API", Version: "1.0"}
	builder := NewBuilder(cfg, reg)
	spec := builder.Build()

	// 检查安全定义是否存在。
	found := false
	for _, s := range spec.Security {
		if s.Name == string(BearerAuth) {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected bearer auth security definition")
	}
}

func TestRegistry_Concurrency(t *testing.T) {
	reg := NewRegistry()
	type EmptyReq struct{}

	// 模拟并发注册（防御性检查）。
	done := make(chan bool)
	for i := 0; i < 50; i++ {
		go func(idx int) {
			rec := &EndpointRecord{
				Method:  "GET",
				Path:    "/api/test",
				ReqType: reflect.TypeOf(EmptyReq{}),
				ResType: reflect.TypeOf(EmptyReq{}),
				Binding: BindingQuery,
			}
			reg.Add(rec)
			reg.Items()
			done <- true
		}(i)
	}

	for i := 0; i < 50; i++ {
		<-done
	}

	items := reg.Items()
	if len(items) != 50 {
		t.Errorf("expected 50 items, got %d", len(items))
	}
}
