package server

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"gin-layout/internal/apidoc"
)

func setupEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestRouteGroup_Protected(t *testing.T) {
	engine := setupEngine()
	reg := apidoc.NewRegistry()
	routes := NewRoutes(engine, WithDocRegistry(reg))

	v1 := routes.Group("/v1").Protected()
	_ = v1.GET("/users", Query[struct{}, struct{}](func(c *gin.Context, req struct{}) (struct{}, error) {
		return struct{}{}, nil
	}))

	records := reg.Items()
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	rec := records[0]
	if rec.GroupDoc.Visibility != apidoc.VisibilityProtected {
		t.Errorf("expected group visibility protected, got %q", rec.GroupDoc.Visibility)
	}
	if len(rec.GroupDoc.Security) != 1 || rec.GroupDoc.Security[0] != apidoc.BearerAuth {
		t.Errorf("expected group security [bearer], got %v", rec.GroupDoc.Security)
	}
}

func TestRouteGroup_Public(t *testing.T) {
	engine := setupEngine()
	reg := apidoc.NewRegistry()
	routes := NewRoutes(engine, WithDocRegistry(reg))

	v1 := routes.Group("/v1").Public()
	_ = v1.GET("/health", Query[struct{}, struct{}](func(c *gin.Context, req struct{}) (struct{}, error) {
		return struct{}{}, nil
	}))

	rec := reg.Items()[0]
	if rec.GroupDoc.Visibility != apidoc.VisibilityPublic {
		t.Errorf("expected group visibility public, got %q", rec.GroupDoc.Visibility)
	}
	if len(rec.GroupDoc.Security) != 0 {
		t.Errorf("expected group security empty, got %v", rec.GroupDoc.Security)
	}
}

func TestRouteGroup_Tag(t *testing.T) {
	engine := setupEngine()
	reg := apidoc.NewRegistry()
	routes := NewRoutes(engine, WithDocRegistry(reg))

	v1 := routes.Group("/v1").Tag("管理")
	users := v1.Group("/users").Tag("用户")
	_ = users.GET("", Query[struct{}, struct{}](func(c *gin.Context, req struct{}) (struct{}, error) {
		return struct{}{}, nil
	}))

	rec := reg.Items()[0]
	want := []string{"管理", "用户"}
	if !reflect.DeepEqual(rec.GroupDoc.Tags, want) {
		t.Errorf("expected tags %v, got %v", want, rec.GroupDoc.Tags)
	}
}

func TestRouteGroup_CloneIsolation(t *testing.T) {
	engine := setupEngine()
	reg := apidoc.NewRegistry()
	routes := NewRoutes(engine, WithDocRegistry(reg))

	v1 := routes.Group("/v1").Protected().Tag("v1")
	users := v1.Group("/users") // inherits v1 defaults
	v1.Tag("extra")               // should not affect users

	_ = users.GET("", Query[struct{}, struct{}](func(c *gin.Context, req struct{}) (struct{}, error) {
		return struct{}{}, nil
	}))

	rec := reg.Items()[0]
	if reflect.DeepEqual(rec.GroupDoc.Tags, []string{"v1", "extra"}) {
		t.Errorf("users group should not inherit v1's later tag, got %v", rec.GroupDoc.Tags)
	}
	want := []string{"v1"}
	if !reflect.DeepEqual(rec.GroupDoc.Tags, want) {
		t.Errorf("expected tags %v, got %v", want, rec.GroupDoc.Tags)
	}
}

func TestRoute_Protected(t *testing.T) {
	engine := setupEngine()
	reg := apidoc.NewRegistry()
	routes := NewRoutes(engine, WithDocRegistry(reg))

	_ = routes.Group("/api").Public().GET("/v1/users", Query[struct{}, struct{}](func(c *gin.Context, req struct{}) (struct{}, error) {
		return struct{}{}, nil
	})).Protected()

	rec := reg.Items()[0]
	if rec.RouteDoc.Visibility == nil || *rec.RouteDoc.Visibility != apidoc.VisibilityProtected {
		t.Errorf("expected route visibility protected")
	}
	if len(rec.RouteDoc.Security) != 1 || rec.RouteDoc.Security[0] != apidoc.BearerAuth {
		t.Errorf("expected route security [bearer], got %v", rec.RouteDoc.Security)
	}
}

func TestRouteGroup_TagDoesNotShareBackingArray(t *testing.T) {
	engine := setupEngine()
	reg := apidoc.NewRegistry()
	routes := NewRoutes(engine, WithDocRegistry(reg))

	v1 := routes.Group("/v1").Tag("a")
	_ = v1.GET("/items", Query[struct{}, struct{}](func(c *gin.Context, req struct{}) (struct{}, error) {
		return struct{}{}, nil
	}))

	users := v1.Group("/users").Tag("b")
	for i := 0; i < 10; i++ {
		users = users.Tag("x")
	}
	_ = users.GET("", Query[struct{}, struct{}](func(c *gin.Context, req struct{}) (struct{}, error) {
		return struct{}{}, nil
	}))

	var v1Tags, usersTags []string
	for _, rec := range reg.Items() {
		if rec.Path == "/v1/items" {
			v1Tags = rec.GroupDoc.Tags
		} else if rec.Path == "/v1/users" {
			usersTags = rec.GroupDoc.Tags
		}
	}

	if !reflect.DeepEqual(v1Tags, []string{"a"}) {
		t.Errorf("parent tags corrupted: got %v", v1Tags)
	}
	wantUsers := append([]string{"a", "b"}, make([]string, 10)...)
	for i := 2; i < len(wantUsers); i++ {
		wantUsers[i] = "x"
	}
	if !reflect.DeepEqual(usersTags, wantUsers) {
		t.Errorf("users tags unexpected: got %v", usersTags)
	}
}
