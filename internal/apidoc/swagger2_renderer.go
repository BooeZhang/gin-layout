package apidoc

import (
	"encoding/json"
	"sort"
	"strings"
)

// ---- Swagger 2.0 传输类型 ----
// 这些类型镜像了 internal/swagger/types.go 中的类型，但保持自包含
// 以避免导入循环（swagger → server → apidoc）。

type swagSpec struct {
	Swagger             string                     `json:"swagger"`
	Info                swagInfo                   `json:"info"`
	Host                string                     `json:"host,omitempty"`
	BasePath            string                     `json:"basePath,omitempty"`
	Schemes             []string                   `json:"schemes,omitempty"`
	Tags                []swagTag                  `json:"tags,omitempty"`
	Paths               map[string]swagPathItem    `json:"paths"`
	Definitions         map[string]swagProperty    `json:"definitions,omitempty"`
	SecurityDefinitions map[string]swagSecurityDef `json:"securityDefinitions,omitempty"`
}

type swagInfo struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

type swagTag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type swagPathItem map[string]swagOperation

type swagOperation struct {
	Summary     string                  `json:"summary,omitempty"`
	Description string                  `json:"description,omitempty"`
	Tags        []string                `json:"tags,omitempty"`
	Parameters  []swagParameter         `json:"parameters,omitempty"`
	Responses   map[string]swagResponse `json:"responses"`
	Deprecated  bool                    `json:"deprecated,omitempty"`
	Security    []map[string][]string   `json:"security,omitempty"`
	Consumes    []string                `json:"consumes,omitempty"`
	Produces    []string                `json:"produces,omitempty"`
}

type swagParameter struct {
	In          string        `json:"in"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Required    bool          `json:"required,omitempty"`
	Type        string        `json:"type,omitempty"`
	Format      string        `json:"format,omitempty"`
	Schema      *swagProperty `json:"schema,omitempty"`
}

type swagResponse struct {
	Description string        `json:"description"`
	Schema      *swagProperty `json:"schema,omitempty"`
}

type swagProperty struct {
	Ref         string                  `json:"$ref,omitempty"`
	Type        string                  `json:"type,omitempty"`
	Format      string                  `json:"format,omitempty"`
	Description string                  `json:"description,omitempty"`
	Properties  map[string]swagProperty `json:"properties,omitempty"`
	Items       *swagProperty           `json:"items,omitempty"`
	Required    []string                `json:"required,omitempty"`
}

type swagSecurityDef struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
	In   string `json:"in,omitempty"`
}

// ---- 渲染器 ----

// Swagger2Renderer 将 SpecModel 转换为 Swagger 2.0 JSON。
type Swagger2Renderer struct{}

// Render 将 SpecModel 转换为 Swagger 2.0 JSON 字节。
func (r *Swagger2Renderer) Render(spec *SpecModel) ([]byte, error) {
	sw := &swagSpec{
		Swagger: "2.0",
		Info: swagInfo{
			Title:       spec.Title,
			Description: spec.Description,
			Version:     spec.Version,
		},
		Host:     spec.Host,
		BasePath: spec.BasePath,
		Schemes:  []string{"http", "https"},
		Paths:    make(map[string]swagPathItem),
	}

	// 标签。
	for _, t := range spec.Tags {
		sw.Tags = append(sw.Tags, swagTag{Name: t.Name, Description: t.Description})
	}

	// 路径。
	for path, pathModel := range spec.Paths {
		item := make(swagPathItem)
		for method, op := range pathModel {
			item[strings.ToLower(method)] = r.convertOperation(op, spec)
		}
		sw.Paths[path] = item
	}

	// 定义。
	if len(spec.Definitions) > 0 {
		sw.Definitions = make(map[string]swagProperty)
		for name, s := range spec.Definitions {
			sw.Definitions[name] = r.convertSchema(s)
		}
	}

	// 安全定义。
	if len(spec.Security) > 0 {
		sw.SecurityDefinitions = make(map[string]swagSecurityDef)
		for _, sec := range spec.Security {
			sw.SecurityDefinitions[sec.Name] = swagSecurityDef{
				Type: "apiKey",
				Name: "Authorization",
				In:   sec.In,
			}
		}
	}

	return json.MarshalIndent(sw, "", "  ")
}

func (r *Swagger2Renderer) convertOperation(op OperationModel, spec *SpecModel) swagOperation {
	swOp := swagOperation{
		Summary:     op.Summary,
		Description: op.Description,
		Tags:        op.Tags,
		Deprecated:  op.Deprecated,
		Responses:   make(map[string]swagResponse),
	}

	// 参数。
	for _, p := range op.Parameters {
		swParam := swagParameter{
			In:          p.In,
			Name:        p.Name,
			Description: p.Description,
			Required:    p.Required,
			Type:        p.Type,
			Format:      p.Format,
		}
		if p.SchemaRef != "" {
			swParam.Schema = &swagProperty{Ref: p.SchemaRef}
		}
		swOp.Parameters = append(swOp.Parameters, swParam)
	}

	// 响应（排序以确保确定性输出）。
	keys := make([]string, 0, len(op.Responses))
	for k := range op.Responses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, status := range keys {
		resp := op.Responses[status]
		swResp := swagResponse{
			Description: resp.Description,
		}
		if resp.SchemaRef != "" {
			swResp.Schema = &swagProperty{Ref: resp.SchemaRef}
		}
		swOp.Responses[status] = swResp
	}

	// 安全。
	for _, s := range op.Security {
		swOp.Security = append(swOp.Security, map[string][]string{
			s.Name: {},
		})
	}

	// Consumes / Produces。
	if r.hasBodyParam(op.Parameters) {
		swOp.Consumes = []string{"application/json"}
	}
	swOp.Produces = []string{"application/json"}

	return swOp
}

func (r *Swagger2Renderer) hasBodyParam(params []ParamModel) bool {
	for _, p := range params {
		if p.In == "body" {
			return true
		}
	}
	return false
}

// 递归地将 SchemaModel 转换为 swagProperty。
func (r *Swagger2Renderer) convertSchema(s SchemaModel) swagProperty {
	prop := swagProperty{
		Type:        s.Type,
		Format:      s.Format,
		Description: s.Description,
	}
	if s.Ref != "" {
		prop = swagProperty{Ref: s.Ref}
		return prop
	}
	if len(s.Properties) > 0 {
		prop.Properties = make(map[string]swagProperty)
		for name, ps := range s.Properties {
			prop.Properties[name] = r.convertSchema(ps)
		}
	}
	if len(s.Required) > 0 {
		prop.Required = s.Required
	}
	if s.Items != nil {
		itemProp := r.convertSchema(*s.Items)
		prop.Items = &itemProp
	}
	return prop
}
