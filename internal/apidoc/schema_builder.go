package apidoc

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

// fullPkgPathRe matches full Go package paths (containing at least one "/")
// followed by a dot — e.g. "gin-layout/internal/user." in generic type names.
var fullPkgPathRe = regexp.MustCompile(`[\w.*-]+(?:/[\w.-]+)+\.`)

// 通过反射从 Go 类型构建 SchemaModel 定义。
type schemaBuilder struct {
	defs map[string]SchemaModel
}

func newSchemaBuilder() *schemaBuilder {
	return &schemaBuilder{defs: make(map[string]SchemaModel)}
}

// 为类型生成唯一名称，例如 "role.CreateRoleReq"。
// 对于泛型类型，会缩短类型参数中的完整模块路径以避免包含 "/"（会破坏 Swagger $ref）。
func typeSchemaName(t reflect.Type) string {
	if t == nil {
		return ""
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	pkg := t.PkgPath()
	name := t.Name()

	if pkg != "" {
		parts := strings.Split(pkg, "/")
		if len(parts) >= 2 {
			pkg = strings.Join(parts[len(parts)-2:], ".")
		}
		return pkg + "." + shortenGenericTypeArgs(name)
	}
	return shortenGenericTypeArgs(name)
}

// 缩短泛型类型参数中出现的完整 Go 模块路径。
// 例如 "PageResult[gin-layout/internal/user.UserItem]" → "PageResult[internal.user.UserItem]"
func shortenGenericTypeArgs(name string) string {
	if !strings.ContainsRune(name, '[') {
		return name
	}
	return fullPkgPathRe.ReplaceAllStringFunc(name, func(match string) string {
		trimmed := strings.TrimSuffix(match, ".")
		var ptrPrefix string
		if strings.HasPrefix(trimmed, "*") {
			ptrPrefix = "*"
			trimmed = trimmed[1:]
		}
		parts := strings.Split(trimmed, "/")
		if len(parts) >= 2 {
			return ptrPrefix + strings.Join(parts[len(parts)-2:], ".") + "."
		}
		return match
	})
}

// 递归地为类型构建 SchemaModel 并注册。
// 返回定义键名，对于基础类型返回空字符串。
func (b *schemaBuilder) addDefinition(t reflect.Type) string {
	if t == nil {
		return ""
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	kind := t.Kind()
	switch {
	case isBasicKind(kind):
		return ""
	case t == reflect.TypeOf(time.Time{}):
		return ""
	case kind == reflect.Slice || kind == reflect.Array:
		elem := t.Elem()
		if isBasicKind(elem.Kind()) {
			return ""
		}
		b.addDefinition(elem)
		return typeSchemaName(t)
	case kind == reflect.Map:
		return ""
	}

	name := typeSchemaName(t)
	if name == "" {
		return ""
	}

	if _, ok := b.defs[name]; ok {
		return name
	}

	schema := SchemaModel{
		Type:       "object",
		Properties: make(map[string]SchemaModel),
	}
	b.defs[name] = schema // 循环引用的占位符

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		if f.Anonymous {
			b.flattenEmbeddedField(name, f)
			continue
		}

		jsonName := fieldName(f)
		if jsonName == "" || jsonName == "-" {
			continue
		}

		fieldSchema := b.fieldToSchema(f)
		fieldSchema.Description = f.Tag.Get("desc")

		binding := f.Tag.Get("binding")
		if strings.Contains(binding, "required") {
			schema.Required = append(schema.Required, jsonName)
		}

		schema.Properties[jsonName] = fieldSchema
	}

	b.defs[name] = schema
	return name
}

// 将嵌入结构体字段合并到父定义中。
func (b *schemaBuilder) flattenEmbeddedField(parentName string, f reflect.StructField) {
	t := f.Type
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	if t == reflect.TypeOf(time.Time{}) {
		return
	}

	embName := typeSchemaName(t)
	if embName != "" {
		if _, ok := b.defs[embName]; !ok {
			b.defs[embName] = SchemaModel{Type: "object", Properties: make(map[string]SchemaModel)}
		}
	}

	parent := b.defs[parentName]
	for i := 0; i < t.NumField(); i++ {
		ef := t.Field(i)
		if !ef.IsExported() {
			continue
		}
		if ef.Anonymous {
			b.flattenEmbeddedField(parentName, ef)
			continue
		}
		jsonName := fieldName(ef)
		if jsonName == "" || jsonName == "-" {
			continue
		}
		fieldSchema := b.fieldToSchema(ef)
		fieldSchema.Description = ef.Tag.Get("desc")
		parent.Properties[jsonName] = fieldSchema

		binding := ef.Tag.Get("binding")
		if strings.Contains(binding, "required") {
			parent.Required = append(parent.Required, jsonName)
		}
	}
	b.defs[parentName] = parent
}

// 将结构体字段映射为 SchemaModel。
func (b *schemaBuilder) fieldToSchema(f reflect.StructField) SchemaModel {
	t := f.Type
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t == reflect.TypeOf(time.Time{}) {
		return SchemaModel{Type: "string", Format: "date-time"}
	}

	switch t.Kind() {
	case reflect.String:
		return SchemaModel{Type: "string"}
	case reflect.Bool:
		return SchemaModel{Type: "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return SchemaModel{Type: "integer", Format: "int32"}
	case reflect.Int64:
		return SchemaModel{Type: "integer", Format: "int64"}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return SchemaModel{Type: "integer", Format: "int64"}
	case reflect.Float32:
		return SchemaModel{Type: "number", Format: "float"}
	case reflect.Float64:
		return SchemaModel{Type: "number", Format: "double"}
	case reflect.Slice, reflect.Array:
		elem := t.Elem()
		if isBasicKind(elem.Kind()) {
			return SchemaModel{Type: "array", Items: &SchemaModel{Type: basicTypeName(elem)}}
		}
		refName := b.addDefinition(elem)
		if refName != "" {
			return SchemaModel{Type: "array", Items: &SchemaModel{Ref: "#/definitions/" + refName}}
		}
		return SchemaModel{Type: "array"}
	case reflect.Map:
		return SchemaModel{Type: "object"}
	case reflect.Struct:
		refName := b.addDefinition(t)
		if refName != "" {
			return SchemaModel{Ref: "#/definitions/" + refName}
		}
		return SchemaModel{Type: "object"}
	case reflect.Interface:
		return SchemaModel{Type: "object"}
	default:
		return SchemaModel{Type: "string"}
	}
}

// 从结构体字段中提取 JSON/字段名称。
func fieldName(f reflect.StructField) string {
	tag := f.Tag.Get("form")
	if tag == "" {
		tag = f.Tag.Get("json")
	}
	if tag == "" {
		return strings.ToLower(f.Name[:1]) + f.Name[1:]
	}
	parts := strings.Split(tag, ",")
	if len(parts) == 0 || parts[0] == "-" {
		return ""
	}
	return parts[0]
}

// 对标量 Go 类型返回 true。
func isBasicKind(k reflect.Kind) bool {
	switch k {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

// 返回基础类型的 JSON Schema 类型名称。
func basicTypeName(t reflect.Type) string {
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

// 如果类型是零字段的结构体则返回 true。
func isEmptyStructType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct && t.NumField() == 0
}
