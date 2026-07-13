package apidoc

import (
	"slices"
	"strings"
)

// 根据 HTTP 方法和路径生成回退摘要
func defaultSummary(method, path string) string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	resource := "resource"
	for _, s := range slices.Backward(segments) {
		if s != "" && !strings.HasPrefix(s, ":") && !strings.HasPrefix(s, "{") {
			resource = s
		}
	}

	switch strings.ToUpper(method) {
	case "POST":
		return "Create " + resource
	case "GET":
		if strings.Contains(path, ":id") || strings.Contains(path, "{id}") {
			return "Get " + resource
		}
		return "List " + resource
	case "PUT", "PATCH":
		return "Update " + resource
	case "DELETE":
		return "Delete " + resource
	default:
		return strings.ToUpper(method) + " " + resource
	}
}

// 从 gin 路由路径中提取 :param 风格的参数
func parsePathParams(path string) []ParamModel {
	var params []ParamModel
	segments := strings.SplitSeq(path, "/")
	for seg := range segments {
		if strings.HasPrefix(seg, ":") {
			name := seg[1:]
			params = append(params, ParamModel{
				In:       "path",
				Name:     name,
				Required: true,
				Type:     "string",
			})
		}
	}
	return params
}

// 将 gin 风格的 /:param 转换为 Swagger 风格的 /{param}
func normalizePath(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if strings.HasPrefix(p, ":") {
			parts[i] = "{" + p[1:] + "}"
		}
	}
	result := strings.Join(parts, "/")
	if !strings.HasPrefix(result, "/") {
		result = "/" + result
	}
	return result
}


