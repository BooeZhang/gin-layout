package apidoc

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gin-layout/config"
)

// Builder 根据注册的端点和配置构造 SpecModel。
type Builder struct {
	cfg     config.APIDocConfig
	reg     *Registry
	schemas *schemaBuilder
}

// NewBuilder 创建一个新的 Builder。
func NewBuilder(cfg config.APIDocConfig, reg *Registry) *Builder {
	return &Builder{
		cfg:     cfg,
		reg:     reg,
		schemas: newSchemaBuilder(),
	}
}

// Build 从所有注册的端点构造 SpecModel。
func (b *Builder) Build() *SpecModel {
	spec := &SpecModel{
		Title:       b.cfg.Title,
		Description: b.cfg.Description,
		Version:     b.cfg.Version,
		Host:        b.cfg.Host,
		BasePath:    b.cfg.BasePath,
		Paths:       make(map[string]PathModel),
		Definitions: make(map[string]SchemaModel),
	}

	if spec.Title == "" {
		spec.Title = DefaultTitle
	}
	if spec.Version == "" {
		spec.Version = DefaultVersion
	}

	// 预注册通用响应信封。
	b.registerResponseEnvelope(spec)

	tagSet := make(map[string]bool)

	for _, rec := range b.reg.Items() {
		if rec == nil {
			continue
		}
		// 检查路由级别的隐藏覆盖。
		if rec.RouteDoc.Hidden != nil && *rec.RouteDoc.Hidden {
			continue
		}
		// 检查组级别的隐藏。
		if rec.RouteDoc.Hidden == nil && rec.GroupDoc.Hidden {
			continue
		}

		op := b.buildOperation(rec, spec)

		swPath := normalizePath(rec.Path)
		if spec.Paths[swPath] == nil {
			spec.Paths[swPath] = make(PathModel)
		}
		method := strings.ToLower(rec.Method)
		spec.Paths[swPath][method] = op

		for _, t := range op.Tags {
			tagSet[t] = true
		}
	}

	// 将 schema 定义复制到 spec 中。
	for name, s := range b.schemas.defs {
		if _, ok := spec.Definitions[name]; !ok {
			spec.Definitions[name] = s
		}
	}

	// 构建排序后的标签列表。
	for t := range tagSet {
		spec.Tags = append(spec.Tags, TagModel{Name: t})
	}
	sort.Slice(spec.Tags, func(i, j int) bool {
		return spec.Tags[i].Name < spec.Tags[j].Name
	})

	// 如果有路由使用 bearer 认证，则添加安全定义。
	for _, rec := range b.reg.Items() {
		secs := b.resolveSecurity(rec)
		for _, s := range secs {
			if s.Name == string(BearerAuth) {
				found := false
				for _, existing := range spec.Security {
					if existing.Name == string(BearerAuth) {
						found = true
						break
					}
				}
				if !found {
					spec.Security = append(spec.Security, SecuritySchemeModel{
						Name: string(BearerAuth),
						Type: "apiKey",
						In:   "header",
					})
				}
			}
		}
	}

	return spec
}

// 为单个端点构建 OperationModel。
func (b *Builder) buildOperation(rec *EndpointRecord, spec *SpecModel) OperationModel {
	op := OperationModel{
		Responses: make(map[string]ResponseModel),
	}

	// 摘要
	op.Summary = rec.RouteDoc.Summary
	if op.Summary == "" {
		op.Summary = defaultSummary(rec.Method, rec.Path)
	}
	op.Description = rec.RouteDoc.Description

	// 标签
	tags := b.resolveTags(rec)
	op.Tags = tags

	// 弃用标记
	if rec.RouteDoc.Deprecated != nil {
		op.Deprecated = *rec.RouteDoc.Deprecated
	} else {
		op.Deprecated = rec.GroupDoc.Deprecated
	}

	// 参数
	op.Parameters = b.buildParameters(rec, spec)

	// 安全
	op.Security = b.resolveSecurity(rec)

	// 响应
	b.buildResponses(&op, rec, spec)

	return op
}

// 返回路由的标签：路由级别 > 组级别 > 路径回退。
func (b *Builder) resolveTags(rec *EndpointRecord) []string {
	if len(rec.RouteDoc.Tags) > 0 {
		return uniq(rec.RouteDoc.Tags)
	}
	if len(rec.GroupDoc.Tags) > 0 {
		return uniq(rec.GroupDoc.Tags)
	}
	// 从路径段回退。
	return b.fallbackTags(rec.Path)
}

// 从路径段中提取标签名称。
func (b *Builder) fallbackTags(path string) []string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	var tags []string
	for _, seg := range segments {
		seg = strings.TrimPrefix(seg, ":")
		seg = strings.TrimPrefix(seg, "{")
		seg = strings.TrimSuffix(seg, "}")
		if seg != "" && seg != "api" && seg != "v1" {
			tags = append(tags, seg)
		}
	}
	return uniq(tags)
}

// 返回路由的安全方案。
func (b *Builder) resolveSecurity(rec *EndpointRecord) []SecuritySchemeModel {
	vis := rec.RouteDoc.Visibility
	if vis == nil {
		vis = &rec.GroupDoc.Visibility
	}

	if *vis == VisibilityPublic {
		return nil
	}

	// 路由级别的安全配置优先。
	secs := rec.RouteDoc.Security
	if len(secs) == 0 {
		secs = rec.GroupDoc.Security
	}

	var result []SecuritySchemeModel
	for _, s := range secs {
		result = append(result, SecuritySchemeModel{Name: string(s)})
	}
	return result
}

// 从路由构建参数模型。
func (b *Builder) buildParameters(rec *EndpointRecord, spec *SpecModel) []ParamModel {
	var params []ParamModel

	// 1. 从路由路径提取路径参数。
	pathParams := parsePathParams(rec.Path)
	params = append(params, pathParams...)

	// 2. 从请求类型提取查询或请求体参数。
	if rec.ReqType != nil {
		reqType := rec.ReqType
		for reqType.Kind() == reflect.Ptr {
			reqType = reqType.Elem()
		}

		switch rec.Binding {
		case BindingQuery, BindingForm:
			qp := b.buildQueryParams(reqType)
			params = append(params, qp...)
		case BindingJSON:
			bp := b.buildBodyParam(reqType)
			params = append(params, bp...)
		case BindingAuto:
			if strings.ToUpper(rec.Method) == "GET" || strings.ToUpper(rec.Method) == "DELETE" {
				qp := b.buildQueryParams(reqType)
				params = append(params, qp...)
			} else {
				bp := b.buildBodyParam(reqType)
				params = append(params, bp...)
			}
		}
	}

	return params
}

// 从请求结构体类型推断查询参数。
func (b *Builder) buildQueryParams(t reflect.Type) []ParamModel {
	if isEmptyStructType(t) {
		return nil
	}

	var params []ParamModel
	b.buildQueryParamsRecursive(t, &params)
	return params
}

func (b *Builder) buildQueryParamsRecursive(t reflect.Type, params *[]ParamModel) {
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		if f.Anonymous {
			ft := f.Type
			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			b.buildQueryParamsRecursive(ft, params)
			continue
		}

		name := fieldName(f)
		if name == "" || name == "-" {
			continue
		}

		param := ParamModel{
			In:          "query",
			Name:        name,
			Description: f.Tag.Get("desc"),
			Type:        b.resolveParamType(f.Type),
			Required:    strings.Contains(f.Tag.Get("binding"), "required"),
		}
		*params = append(*params, param)
	}
}

// 从请求结构体类型创建请求体参数。
func (b *Builder) buildBodyParam(t reflect.Type) []ParamModel {
	if isEmptyStructType(t) {
		return nil
	}

	refName := b.schemas.addDefinition(t)
	if refName == "" && t.Kind() == reflect.Slice {
		refName = b.schemas.addDefinition(t)
	}
	if refName == "" {
		return nil
	}

	return []ParamModel{{
		In:        "body",
		Name:      "body",
		Required:  true,
		SchemaRef: "#/definitions/" + refName,
	}}
}

// 将 reflect.Type 映射为查询/路径参数的 Swagger 类型字符串。
func (b *Builder) resolveParamType(t reflect.Type) string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	default:
		return "string"
	}
}

// 构建操作的响应和错误响应。
func (b *Builder) buildResponses(op *OperationModel, rec *EndpointRecord, spec *SpecModel) {
	// 成功响应 (200)。
	resSchema := b.buildResponseSchema(rec, spec)
	op.Responses["200"] = ResponseModel{
		Description: "成功",
		SchemaRef:   resSchema,
	}

	// 错误响应。
	errors := b.mergeErrors(rec)
	for _, e := range errors {
		status := fmt.Sprintf("%d", e.HTTPStatus)
		if _, exists := op.Responses[status]; !exists {
			op.Responses[status] = ResponseModel{
				Description: e.Message,
			}
		}
	}
}

// 生成带有真实数据类型的响应信封 schema。
// 返回类似 "#/definitions/common.Response[user.UserItem]" 的定义键。
func (b *Builder) buildResponseSchema(rec *EndpointRecord, spec *SpecModel) string {
	if rec.ResType == nil {
		return "#/definitions/common.Response"
	}

	resType := rec.ResType
	for resType.Kind() == reflect.Ptr {
		resType = resType.Elem()
	}

	// 对于空结构体，使用基础的 common.Response。
	if resType.Kind() == reflect.Struct && resType.NumField() == 0 {
		return "#/definitions/common.Response"
	}

	// 将响应类型注册为定义。
	dataRef := b.schemas.addDefinition(resType)
	if dataRef == "" {
		return "#/definitions/common.Response"
	}

	// 创建带有具体数据类型的包装响应定义。
	wrapperName := "common.Response[" + dataRef + "]"
	if _, ok := spec.Definitions[wrapperName]; !ok {
		spec.Definitions[wrapperName] = SchemaModel{
			Type: "object",
			Properties: map[string]SchemaModel{
				"code":    {Type: "integer", Description: "业务状态码，0 表示成功"},
				"message": {Type: "string", Description: "提示信息"},
				"data":    {Ref: "#/definitions/" + dataRef},
			},
		}
	}

	return "#/definitions/" + wrapperName
}

// 预注册基础的 common.Response 定义。
func (b *Builder) registerResponseEnvelope(spec *SpecModel) {
	spec.Definitions["common.Response"] = SchemaModel{
		Type: "object",
		Properties: map[string]SchemaModel{
			"code":    {Type: "integer", Description: "业务状态码，0 表示成功"},
			"message": {Type: "string", Description: "提示信息"},
			"data":    {Type: "object", Description: "响应数据"},
		},
	}
}

// 合并组级别的默认错误与路由级别的覆盖。
// 相同 HTTP 状态码的路由错误会覆盖组错误。
func (b *Builder) mergeErrors(rec *EndpointRecord) []ErrorSpec {
	seen := make(map[int]bool)
	var result []ErrorSpec

	// 路由错误优先。
	for _, e := range rec.RouteDoc.Errors {
		if !seen[e.HTTPStatus] {
			result = append(result, e)
			seen[e.HTTPStatus] = true
		}
	}

	// 组默认错误填补空缺。
	for _, e := range rec.GroupDoc.DefaultErrors {
		if !seen[e.HTTPStatus] {
			result = append(result, e)
			seen[e.HTTPStatus] = true
		}
	}

	// 如果完全没有定义错误，则添加合理的默认值。
	if len(result) == 0 {
		result = append(result, ErrorSpec{HTTPStatus: 400, Message: "请求参数错误"})
		result = append(result, ErrorSpec{HTTPStatus: 500, Message: "服务器内部错误"})

		// 为非公开路由添加认证错误。
		hasSecurity := len(rec.RouteDoc.Security) > 0 || len(rec.GroupDoc.Security) > 0
		if hasSecurity {
			result = append(result, ErrorSpec{HTTPStatus: 401, Message: "未授权"})
			result = append(result, ErrorSpec{HTTPStatus: 403, Message: "无权限"})
		}
	}

	return result
}

// 对字符串切片去重并保持顺序。
func uniq(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
