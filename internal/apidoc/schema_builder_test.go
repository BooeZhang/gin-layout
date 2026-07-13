package apidoc

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTypeSchemaName(t *testing.T) {
	tests := []struct {
		name string
		t    reflect.Type
		want string
	}{
		{"nil", nil, ""},
		{"time.Time", reflect.TypeOf(time.Time{}), "time.Time"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typeSchemaName(tt.t)
			if got != tt.want {
				t.Errorf("typeSchemaName(%v) = %q, want %q", tt.t, got, tt.want)
			}
		})
	}
}

func TestSchemaBuilder_BasicTypes(t *testing.T) {
	b := newSchemaBuilder()

	type SimpleReq struct {
		Name  string `json:"name" desc:"用户名称" binding:"required"`
		Age   int    `json:"age" desc:"年龄"`
		Email string `json:"email"`
	}

	name := b.addDefinition(reflect.TypeOf(SimpleReq{}))
	t.Logf("definition name: %s", name)
	if name == "" {
		t.Fatal("expected non-empty definition name")
	}

	def, ok := b.defs[name]
	if !ok {
		t.Fatalf("definition %q not found in defs", name)
	}

	if def.Type != "object" {
		t.Errorf("expected type=object, got %q", def.Type)
	}

	nameProp, ok := def.Properties["name"]
	if !ok {
		t.Fatal("expected 'name' property")
	}
	if nameProp.Type != "string" {
		t.Errorf("expected name.type=string, got %q", nameProp.Type)
	}
	if nameProp.Description != "用户名称" {
		t.Errorf("expected name.description='用户名称', got %q", nameProp.Description)
	}

	ageProp, ok := def.Properties["age"]
	if !ok {
		t.Fatal("expected 'age' property")
	}
	if ageProp.Type != "integer" {
		t.Errorf("expected age.type=integer, got %q", ageProp.Type)
	}

	// 检查必填字段
	found := false
	for _, r := range def.Required {
		if r == "name" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'name' to be required")
	}
}

func TestSchemaBuilder_NestedStruct(t *testing.T) {
	b := newSchemaBuilder()

	type Address struct {
		City string `json:"city"`
		Zip  string `json:"zip"`
	}

	type UserReq struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	name := b.addDefinition(reflect.TypeOf(UserReq{}))
	if name == "" {
		t.Fatal("expected non-empty definition name")
	}

	def := b.defs[name]
	addrProp, ok := def.Properties["address"]
	if !ok {
		t.Fatal("expected 'address' property")
	}
	if addrProp.Ref == "" {
		t.Error("expected address to be a $ref to nested definition")
	}
}

func TestSchemaBuilder_EmbeddedStruct(t *testing.T) {
	b := newSchemaBuilder()

	type BaseModel struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	type ProductReq struct {
		BaseModel
		Price float64 `json:"price"`
	}

	name := b.addDefinition(reflect.TypeOf(ProductReq{}))
	if name == "" {
		t.Fatal("expected non-empty definition name")
	}

	def := b.defs[name]
	// 嵌入字段应被展平。
	if _, ok := def.Properties["id"]; !ok {
		t.Error("expected 'id' property from embedded BaseModel")
	}
	if _, ok := def.Properties["name"]; !ok {
		t.Error("expected 'name' property from embedded BaseModel")
	}
	if _, ok := def.Properties["price"]; !ok {
		t.Error("expected 'price' property")
	}
}

func TestSchemaBuilder_Slice(t *testing.T) {
	b := newSchemaBuilder()

	type TagReq struct {
		Labels []string `json:"labels"`
	}

	name := b.addDefinition(reflect.TypeOf(TagReq{}))
	if name == "" {
		t.Fatal("expected non-empty definition name")
	}

	def := b.defs[name]
	labelsProp, ok := def.Properties["labels"]
	if !ok {
		t.Fatal("expected 'labels' property")
	}
	if labelsProp.Type != "array" {
		t.Errorf("expected array type, got %q", labelsProp.Type)
	}
	if labelsProp.Items == nil {
		t.Fatal("expected items for array type")
	}
	if labelsProp.Items.Type != "string" {
		t.Errorf("expected items.type=string, got %q", labelsProp.Items.Type)
	}
}

func TestSchemaBuilder_TimeField(t *testing.T) {
	b := newSchemaBuilder()

	type EventReq struct {
		CreatedAt time.Time `json:"created_at"`
	}

	name := b.addDefinition(reflect.TypeOf(EventReq{}))
	if name == "" {
		t.Fatal("expected non-empty definition name")
	}

	def := b.defs[name]
	createdAt, ok := def.Properties["created_at"]
	if !ok {
		t.Fatal("expected 'created_at' property")
	}
	if createdAt.Type != "string" {
		t.Errorf("expected string type, got %q", createdAt.Type)
	}
	if createdAt.Format != "date-time" {
		t.Errorf("expected date-time format, got %q", createdAt.Format)
	}
}

func TestSchemaBuilder_EmptyStruct(t *testing.T) {
	b := newSchemaBuilder()

	type EmptyReq struct{}
	name := b.addDefinition(reflect.TypeOf(EmptyReq{}))
	if name == "" {
		t.Fatal("expected non-empty definition name for named empty struct")
	}

	def := b.defs[name]
	if def.Type != "object" {
		t.Errorf("expected object type, got %q", def.Type)
	}
	if len(def.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(def.Properties))
	}
}

func TestIsEmptyStructType(t *testing.T) {
	if !isEmptyStructType(reflect.TypeOf(struct{}{})) {
		t.Error("expected empty struct")
	}
	if isEmptyStructType(reflect.TypeOf(struct{ X int }{})) {
		t.Error("expected non-empty struct")
	}
	if !isEmptyStructType(reflect.TypeOf(&struct{}{})) {
		t.Error("expected pointer-to-empty-struct to be empty")
	}
}

func TestIsBasicKind(t *testing.T) {
	if !isBasicKind(reflect.String) {
		t.Error("string should be basic")
	}
	if isBasicKind(reflect.Struct) {
		t.Error("struct should not be basic")
	}
	if isBasicKind(reflect.Slice) {
		t.Error("slice should not be basic")
	}
}

func TestFieldName(t *testing.T) {
	type testStruct struct {
		JSONField  string `json:"json_name"`
		FormField  string `form:"form_name"`
		BothField  string `json:"json_both" form:"form_both"`
		NoTagField string
		SkipField  string `json:"-"`
	}

	rt := reflect.TypeOf(testStruct{})

	tests := []struct {
		idx  int
		want string
	}{
		{0, "json_name"},
		{1, "form_name"},
		{2, "form_both"}, // form takes priority for query params
		{3, "noTagField"},
		{4, ""}, // skipped
	}

	for _, tt := range tests {
		name := fieldName(rt.Field(tt.idx))
		if name != tt.want {
			t.Errorf("fieldName(%s) = %q, want %q", rt.Field(tt.idx).Name, name, tt.want)
		}
	}
}

func TestShortenGenericTypeArgs(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Non-generic name: no change
		{"UserItem", "UserItem"},
		{"goldenLoginReq", "goldenLoginReq"},
		// Generic with full package paths
		{
			"PageResult[gin-layout/internal/user.UserItem]",
			"PageResult[internal.user.UserItem]",
		},
		// Generic with short paths already: no change
		{"PageResult[user.UserItem]", "PageResult[user.UserItem]"},
		// Single-segment path (no slash): no change
		{"PageResult[UserItem]", "PageResult[UserItem]"},
		// Pointer type arg
		{
			"PageResult[*gin-layout/internal/user.UserItem]",
			"PageResult[*internal.user.UserItem]",
		},
		// Multiple type args
		{
			"MapResult[gin-layout/internal/user.UserItem, gin-layout/internal/role.RoleItem]",
			"MapResult[internal.user.UserItem, internal.role.RoleItem]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := shortenGenericTypeArgs(tt.input)
			if got != tt.want {
				t.Errorf("shortenGenericTypeArgs(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTypeSchemaName_GenericPageResult(t *testing.T) {
	// Simulates domain.PageResult[user.UserItem] — generic type with cross-package type arg
	type testPageResult[T any] struct {
		Items []T `json:"items"`
	}

	type testUserItem struct {
		ID int64 `json:"id"`
	}

	resType := reflect.TypeOf(testPageResult[testUserItem]{})
	got := typeSchemaName(resType)

	// The name must NOT contain "/" (which breaks Swagger $ref)
	if strings.Contains(got, "/") {
		t.Errorf("typeSchemaName must not contain '/', got %q", got)
	}

	// The name should contain the shortened package info
	if got == "" {
		t.Error("typeSchemaName must not be empty")
	}
	t.Logf("typeSchemaName: %s", got)
}
