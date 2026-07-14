package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"gin-layout/internal/apidoc"
	"gin-layout/internal/auth"
	"gin-layout/internal/common"
	"gin-layout/internal/health"
	"gin-layout/internal/menu"
	mw "gin-layout/internal/middleware"
	"gin-layout/internal/role"
	"gin-layout/internal/server"
	"gin-layout/internal/user"
)

// AdminRouter 注册所有 API 路由。
type AdminRouter struct {
	authH       *auth.Handler
	healthH     *health.Handler
	userH       *user.Handler
	roleH       *role.Handler
	menuH       *menu.Handler
	docRegistry *apidoc.Registry
	tokens      common.TokenManager
	policy      common.PolicyManager
	permMap     common.PermissionResolver
	log         *zerolog.Logger
}

// AdminRouterConfig 保存 AdminRouter 的依赖项
type AdminRouterConfig struct {
	Auth        *auth.Handler
	Health      *health.Handler
	User        *user.Handler
	Role        *role.Handler
	Menu        *menu.Handler
	DocRegistry *apidoc.Registry
	Tokens      common.TokenManager
	Policy      common.PolicyManager
	PermMap     common.PermissionResolver
	Log         *zerolog.Logger
}

// NewAdminRouter 创建一个新的 AdminRouter 对象
func NewAdminRouter(cfg AdminRouterConfig) *AdminRouter {
	return &AdminRouter{
		authH:       cfg.Auth,
		healthH:     cfg.Health,
		userH:       cfg.User,
		roleH:       cfg.Role,
		menuH:       cfg.Menu,
		docRegistry: cfg.DocRegistry,
		tokens:      cfg.Tokens,
		policy:      cfg.Policy,
		permMap:     cfg.PermMap,
		log:         cfg.Log,
	}
}

// Register implements server.Router.
func (r *AdminRouter) Register(engine *gin.Engine) {
	authMW := mw.Auth(r.tokens)
	rbacMW := mw.RBAC(r.policy, r.permMap, r.log)

	routes := server.NewRoutes(engine, server.WithDocRegistry(r.docRegistry))
	routes.GET("/health", server.Query(r.healthH.Check)).Summary("健康检查").Tag("系统")

	api := routes.Group("/api")

	auth := api.Group("/auth").Public().Tag("认证")
	auth.POST("/login", server.JSON(r.authH.Login)).Summary("登录")
	auth.POST("/refresh-token", server.JSON(r.authH.RefreshToken)).Summary("刷新令牌")
	auth.POST("/logout", server.JSON(r.authH.Logout), authMW).Summary("登出").Protected()

	v1 := api.Group("/v1", authMW, rbacMW).Protected()

	users := v1.Group("/users").Tag("用户")
	users.POST("", server.JSON(r.userH.Create)).Summary("创建用户")
	users.GET("", server.Query(r.userH.List)).Summary("用户列表")
	users.GET("/details", server.Query(r.userH.GetDetails)).Summary("当前用户详情")
	users.PUT("/:id", server.JSON(r.userH.Update)).Summary("更新用户")
	users.DELETE("/:id", server.Query(r.userH.Delete)).Summary("删除用户")
	users.GET("/menus", server.Query(r.userH.GetMenus)).Summary("获取当前用户菜单")

	roles := v1.Group("/roles").Tag("角色")
	roles.POST("", server.JSON(r.roleH.Create)).Summary("创建角色")
	roles.GET("", server.Query(r.roleH.List)).Summary("角色列表")
	roles.GET("/all", server.Query(r.roleH.GetAll)).Summary("所有角色")
	roles.GET("/:id", server.Query(r.roleH.GetOne)).Summary("角色详情")
	roles.PUT("/:id", server.JSON(r.roleH.Update)).Summary("更新角色")
	roles.DELETE("/:id", server.Query(r.roleH.Delete)).Summary("删除角色")
	roles.PUT("/user-add/:id", server.JSON(r.roleH.UserAdd)).Summary("添加用户到角色")
	roles.PUT("/user-remove/:id", server.JSON(r.roleH.UserRemove)).Summary("从角色移除用户")

	menus := v1.Group("/menus").Tag("菜单")
	menus.POST("", server.JSON(r.menuH.Create)).Summary("创建菜单")
	menus.GET("", server.Query(r.menuH.List)).Summary("菜单列表")
	menus.GET("/:id", server.Query(r.menuH.GetOne)).Summary("菜单详情")
	menus.PUT("/:id", server.JSON(r.menuH.Update)).Summary("更新菜单")
	menus.DELETE("/:id", server.Query(r.menuH.Delete)).Summary("删除菜单")
}
