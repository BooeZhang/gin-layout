package apidoc

import (
	"cmp"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"gin-layout/config"
)

//go:embed static/redoc.standalone.js
var redocJS []byte

// Publisher 在运行时提供 API 文档 JSON 和 Redoc UI。
// Spec 惰性构建（首次请求 /doc.json 时），也可通过 Build() 在启动时主动构建以快速失败。
type Publisher struct {
	cfg      config.APIDocConfig
	registry *Registry
	jsonPath string
	uiPath   string
	enabled  bool

	mu       sync.RWMutex
	json     []byte
	built    bool
	buildErr error
}

// NewPublisher 创建一个 Publisher 但不立即构建 spec。
// 注册表在路由注册期间填充；spec 在首次请求或显式调用 Build() 时构建。
func NewPublisher(cfg config.APIDocConfig, reg *Registry) *Publisher {
	if !cfg.Enabled {
		return &Publisher{enabled: false}
	}

	jsonPath := cmp.Or(cfg.JSONPath, DefaultJSONPath)
	uiPath := cmp.Or(cfg.UIPath, DefaultUIPath)

	return &Publisher{
		cfg:      cfg,
		registry: reg,
		jsonPath: jsonPath,
		uiPath:   uiPath,
		enabled:  true,
	}
}

// Build 主动构建 spec。在所有路由注册完成后调用以快速失败。
func (p *Publisher) Build() error {
	_, err := p.buildSpec()
	return err
}

// UIPath 返回去掉通配符后的文档 UI 路径，用于日志打印。例如 "/docs/"。
func (p *Publisher) UIPath() string {
	return strings.TrimSuffix(p.uiPath, "/*any")
}

// Register 实现 server.Router，注册文档 JSON 与 Redoc UI 端点。
func (p *Publisher) Register(engine *gin.Engine) {
	if !p.enabled {
		return
	}
	engine.GET(p.jsonPath, p.handleJSON)
	engine.GET(p.uiPath, p.handleUI)
}

func (p *Publisher) handleJSON(c *gin.Context) {
	data, err := p.buildSpec()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

func (p *Publisher) handleUI(c *gin.Context) {
	// 对 /docs/redoc.standalone.js 这类资源请求返回嵌入的 JS 文件。
	if strings.HasSuffix(c.Param("any"), ".js") {
		c.Data(http.StatusOK, "application/javascript; charset=utf-8", redocJS)
		return
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
body { margin: 0; padding: 0; }
redoc { display: block; }
</style>
</head>
<body>
<redoc spec-url="%s"></redoc>
<script src="redoc.standalone.js"></script>
</body>
</html>
`, htmlEscaper.Replace(p.cfg.Title), htmlEscaper.Replace(p.jsonPath))

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// htmlEscaper 防止 HTML 模板中的标题或 URL 被意外转义。
var htmlEscaper = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	"\"", "&#34;",
	"'", "&#39;",
)

func (p *Publisher) buildSpec() ([]byte, error) {
	p.mu.RLock()
	if p.built {
		data, err := p.json, p.buildErr
		p.mu.RUnlock()
		return data, err
	}
	p.mu.RUnlock()

	p.mu.Lock()
	defer p.mu.Unlock()
	if p.built {
		return p.json, p.buildErr
	}

	builder := NewBuilder(p.cfg, p.registry)
	spec := builder.Build()

	renderer := &Swagger2Renderer{}
	data, err := renderer.Render(spec)
	if err != nil {
		p.buildErr = fmt.Errorf("apidoc: render failed: %w", err)
		p.built = true
		return nil, p.buildErr
	}

	p.json = data
	p.built = true
	return data, nil
}
