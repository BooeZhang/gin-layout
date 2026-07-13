// Package apidoc 提供一次构建多次服务的 API 文档子系统。
// 在路由注册期间，端点及其文档元数据被收集到 Registry 中。
// 启动时，Builder 合并组默认值与路由覆盖，执行 schema 反射，
// 并生成内部的 SpecModel。渲染器将模型转换为 Swagger 2.0 JSON，
// Publisher 在运行时提供不可变的字节数据。
package apidoc

import "reflect"

// BindingMode 描述类型化处理器如何读取请求数据。
type BindingMode string

const (
	BindingAuto  BindingMode = "auto"
	BindingJSON  BindingMode = "json"
	BindingQuery BindingMode = "query"
	BindingForm  BindingMode = "form"
)

// Visibility 描述路由是公开还是受保护。
type Visibility string

const (
	VisibilityPublic    Visibility = "public"
	VisibilityProtected Visibility = "protected"
)

// SecurityScheme 是文档级别的安全标签。
type SecurityScheme string

const (
	BearerAuth SecurityScheme = "bearer"
)

// ErrorSpec 描述 HTTP 错误响应。
type ErrorSpec struct {
	HTTPStatus int
	Message    string
}

// DocDefaults 保存从路由组传播到其所有子路由的文档元数据。
// 路由级别的覆盖优先。
type DocDefaults struct {
	Tags          []string
	Security      []SecurityScheme
	DefaultErrors []ErrorSpec
	Hidden        bool
	Deprecated    bool
	Visibility    Visibility
}

// RouteDoc 保存每个路由的文档元数据。指针字段用于区分"未设置"和"设为 false"。
type RouteDoc struct {
	Summary     string
	Description string
	Tags        []string
	Security    []SecurityScheme
	Errors      []ErrorSpec
	Hidden      *bool
	Deprecated  *bool
	Visibility  *Visibility
}

// EndpointRecord 存储为单个路由构建文档所需的所有元数据。
// 在路由注册期间由 Registry 收集。
type EndpointRecord struct {
	Method   string
	Path     string
	ReqType  reflect.Type
	ResType  reflect.Type
	Binding  BindingMode
	GroupDoc DocDefaults
	RouteDoc RouteDoc
}

// ---- 内部 Spec 模型（与渲染器无关） ----

// SpecModel 是顶层 API 规范模型。
type SpecModel struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Paths       map[string]PathModel // path → PathModel
	Definitions map[string]SchemaModel
	Security    []SecuritySchemeModel
	Tags        []TagModel
}

// TagModel 表示文档标签分组。
type TagModel struct {
	Name        string
	Description string
}

// PathModel 将 HTTP 方法（小写）映射到 OperationModel。
type PathModel map[string]OperationModel

// OperationModel 描述单个 API 操作。
type OperationModel struct {
	Summary     string
	Description string
	Tags        []string
	Parameters  []ParamModel
	Responses   map[string]ResponseModel // status code → response
	Security    []SecuritySchemeModel
	Deprecated  bool
}

// ParamModel 描述操作参数（查询、路径或请求体）。
type ParamModel struct {
	In          string // "query", "path", "body"
	Name        string
	Description string
	Required    bool
	Type        string // for query/path params
	Format      string // for query/path params
	SchemaRef   string // for body params: "#/definitions/Foo"
	ItemsRef    string // for array body params
}

// ResponseModel 描述特定状态码的 API 响应。
type ResponseModel struct {
	Description string
	SchemaRef   string // "#/definitions/common.Response" or custom
}

// SchemaModel 表示 JSON Schema 定义条目。
type SchemaModel struct {
	Ref         string
	Type        string
	Format      string
	Description string
	Properties  map[string]SchemaModel
	Items       *SchemaModel
	Required    []string
}

// SecuritySchemeModel 表示安全定义（例如 bearer 认证）。
type SecuritySchemeModel struct {
	Name string // e.g., "bearer"
	Type string // e.g., "apiKey", "basic"
	In   string // e.g., "header"
}
