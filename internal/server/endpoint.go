package server

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"

	"gin-layout/internal/apidoc"
	"gin-layout/internal/common"
)

// TypedHandler 是带类型请求和响应的 handler 签名。
type TypedHandler[Req, Res any] func(c *gin.Context, req Req) (Res, error)

// WrapJSON 将 TypedHandler 包装为 gin.HandlerFunc，使用 ShouldBindJSON。
func WrapJSON[Req, Res any](h TypedHandler[Req, Res]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Error(err)
			return
		}
		res, err := h(c, req)
		if err != nil {
			c.Error(err)
			return
		}
		common.OK(c, res)
	}
}

// WrapQuery 将 TypedHandler 包装为 gin.HandlerFunc，使用 ShouldBindQuery。
func WrapQuery[Req, Res any](h TypedHandler[Req, Res]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(err)
			return
		}
		res, err := h(c, req)
		if err != nil {
			c.Error(err)
			return
		}
		common.OK(c, res)
	}
}

// WrapForm 将 TypedHandler 包装为 gin.HandlerFunc，使用 ShouldBind（混合绑定）。
func WrapForm[Req, Res any](h TypedHandler[Req, Res]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		if err := c.ShouldBind(&req); err != nil {
			c.Error(err)
			return
		}
		res, err := h(c, req)
		if err != nil {
			c.Error(err)
			return
		}
		common.OK(c, res)
	}
}

// WrapAuto 根据 HTTP 方法自动选择绑定方式：
//   - GET / DELETE → ShouldBindQuery
//   - 其他 → ShouldBindJSON
func WrapAuto[Req, Res any](h TypedHandler[Req, Res]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Req
		var err error
		switch c.Request.Method {
		case http.MethodGet, http.MethodDelete:
			err = c.ShouldBindQuery(&req)
		default:
			err = c.ShouldBindJSON(&req)
		}
		if err != nil {
			c.Error(err)
			return
		}
		res, err := h(c, req)
		if err != nil {
			c.Error(err)
			return
		}
		common.OK(c, res)
	}
}

// RouteHandler 是一个 gin 处理器，加上文档所需的类型元数据。
type RouteHandler struct {
	handler gin.HandlerFunc
	reqType reflect.Type
	resType reflect.Type
	binding apidoc.BindingMode
}

// JSON 绑定
func JSON[Req, Res any](h TypedHandler[Req, Res]) RouteHandler {
	return RouteHandler{
		handler: WrapJSON(h),
		reqType: typeOf[Req](),
		resType: typeOf[Res](),
		binding: apidoc.BindingJSON,
	}
}

// Query 绑定
func Query[Req, Res any](h TypedHandler[Req, Res]) RouteHandler {
	return RouteHandler{
		handler: WrapQuery(h),
		reqType: typeOf[Req](),
		resType: typeOf[Res](),
		binding: apidoc.BindingQuery,
	}
}

// Form 混合绑定
func Form[Req, Res any](h TypedHandler[Req, Res]) RouteHandler {
	return RouteHandler{
		handler: WrapForm(h),
		reqType: typeOf[Req](),
		resType: typeOf[Res](),
		binding: apidoc.BindingForm,
	}
}

// Auto 根据请求方法选择不同绑定器
func Auto[Req, Res any](h TypedHandler[Req, Res]) RouteHandler {
	return RouteHandler{
		handler: WrapAuto(h),
		reqType: typeOf[Req](),
		resType: typeOf[Res](),
		binding: apidoc.BindingAuto,
	}
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// RouteGroup 在 Gin group 上添加文档元数据
type RouteGroup struct {
	group       *gin.RouterGroup
	docDefaults *apidoc.DocDefaults
	docRegistry *apidoc.Registry
}

// RouteGroupOption 配置一个 RouteGroup
type RouteGroupOption func(*RouteGroup)

// WithDocRegistry 设置文档存储
func WithDocRegistry(reg *apidoc.Registry) RouteGroupOption {
	return func(g *RouteGroup) {
		g.docRegistry = reg
	}
}

// NewRoutes 创建根路由组
func NewRoutes(engine *gin.Engine, opts ...RouteGroupOption) *RouteGroup {
	routes := &RouteGroup{group: &engine.RouterGroup}
	for _, opt := range opts {
		opt(routes)
	}
	return routes
}

// Group 创建一个子路由组
func (g *RouteGroup) Group(relativePath string, middleware ...gin.HandlerFunc) *RouteGroup {
	return &RouteGroup{
		group:       g.group.Group(relativePath, middleware...),
		docDefaults: g.docDefaults,
		docRegistry: g.docRegistry,
	}
}

// Doc 为此 group 设置文档默认值，这些默认值会传递到所有子路由和嵌套组。路由级别的覆盖优先。
func (g *RouteGroup) Doc(defaults apidoc.DocDefaults) *RouteGroup {
	g.docDefaults = &defaults
	return g
}

// Use 在此 group 上使用中间件，用于之后注册的路由
func (g *RouteGroup) Use(middleware ...gin.HandlerFunc) *RouteGroup {
	g.group.Use(middleware...)
	return g
}

func (g *RouteGroup) GET(path string, h RouteHandler, middleware ...gin.HandlerFunc) *Route {
	return g.handle(http.MethodGet, path, h, middleware...)
}

func (g *RouteGroup) POST(path string, h RouteHandler, middleware ...gin.HandlerFunc) *Route {
	return g.handle(http.MethodPost, path, h, middleware...)
}

func (g *RouteGroup) PUT(path string, h RouteHandler, middleware ...gin.HandlerFunc) *Route {
	return g.handle(http.MethodPut, path, h, middleware...)
}

func (g *RouteGroup) PATCH(path string, h RouteHandler, middleware ...gin.HandlerFunc) *Route {
	return g.handle(http.MethodPatch, path, h, middleware...)
}

func (g *RouteGroup) DELETE(path string, h RouteHandler, middleware ...gin.HandlerFunc) *Route {
	return g.handle(http.MethodDelete, path, h, middleware...)
}

func (g *RouteGroup) handle(method, path string, h RouteHandler, middleware ...gin.HandlerFunc) *Route {
	handlers := make([]gin.HandlerFunc, 0, len(middleware)+1)
	handlers = append(handlers, middleware...)
	handlers = append(handlers, h.handler)
	g.group.Handle(method, path, handlers...)

	fullPath := joinPaths(g.group.BasePath(), path)

	route := &Route{}

	// 立即将其推送到 apidoc 注册表中，并内置组默认值。
	// 记录是通过指针存储的，因此随后的 Route 方法调用（Summary、Tag 等）会在原地修改它，注册表将看到最终状态。
	if g.docRegistry != nil {
		route.record = &apidoc.EndpointRecord{
			Method:  method,
			Path:    fullPath,
			ReqType: h.reqType,
			ResType: h.resType,
			Binding: h.binding,
		}
		if g.docDefaults != nil {
			route.record.GroupDoc = *g.docDefaults
		}
		g.docRegistry.Add(route.record)
	}

	return route
}

func joinPaths(base, relative string) string {
	base = strings.TrimRight(base, "/")
	if base == "" {
		base = "/"
	}
	if relative == "" || relative == "/" {
		return base
	}
	if strings.HasPrefix(relative, "/") {
		if base == "/" {
			return relative
		}
		return base + relative
	}
	if base == "/" {
		return "/" + relative
	}
	return base + "/" + relative
}

// Route 为已注册的路由可链式调用的元数据设置
type Route struct {
	record *apidoc.EndpointRecord
}

// Summary 设置路由简短简介
func (r *Route) Summary(text string) *Route {
	if r.record != nil {
		r.record.RouteDoc.Summary = text
	}
	return r
}

// Desc 设置路由完整介绍
func (r *Route) Desc(text string) *Route {
	if r.record != nil {
		r.record.RouteDoc.Description = text
	}
	return r
}

// Tag 为路由添加一个或多个标签
func (r *Route) Tag(tags ...string) *Route {
	if r.record != nil {
		r.record.RouteDoc.Tags = append(r.record.RouteDoc.Tags, tags...)
	}
	return r
}

// Error 为此路由声明错误响应。附加到任何组级默认错误
func (r *Route) Error(httpStatus int, message string) *Route {
	if r.record != nil {
		r.record.RouteDoc.Errors = append(r.record.RouteDoc.Errors, apidoc.ErrorSpec{
			HTTPStatus: httpStatus,
			Message:    message,
		})
	}
	return r
}

// Public 标记该路由为公开路由无需身份认证
func (r *Route) Public() *Route {
	if r.record != nil {
		v := apidoc.VisibilityPublic
		r.record.RouteDoc.Visibility = &v
		r.record.RouteDoc.Security = []apidoc.SecurityScheme{}
	}
	return r
}

// Hide 隐藏路由不在文档中显示
func (r *Route) Hide() *Route {
	if r.record != nil {
		hidden := true
		r.record.RouteDoc.Hidden = &hidden
	}
	return r
}

// Deprecated 文档中将该路线标记为已弃用。
func (r *Route) Deprecated() *Route {
	if r.record != nil {
		dep := true
		r.record.RouteDoc.Deprecated = &dep
	}
	return r
}
