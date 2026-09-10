package bootstrap

import (
	"context"
	"fmt"
	"gin-layout/internal/infra/log"
	"time"

	"gin-layout/config"
	"gin-layout/internal/bootstrap/initializer"
	"gin-layout/internal/domain/health"
	"gin-layout/internal/domain/menu"
	"gin-layout/internal/domain/role"
	"gin-layout/internal/domain/sysuser"
	"gin-layout/internal/http/apidoc"
	"gin-layout/internal/http/router"
	"gin-layout/internal/http/server"
	"gin-layout/internal/infra"
	"gin-layout/internal/shared/common"
	"gin-layout/internal/shared/policy"
	"gin-layout/internal/shared/token"

	"github.com/gin-gonic/gin"
)

const blacklistCleanupInterval = 1 * time.Hour

type App struct {
	HTTPServer    *server.Server
	DB            *infra.Database
	Redis         *infra.RedisClient
	cleanupCancel context.CancelFunc
}

type appInfra struct {
	db    *infra.Database
	redis *infra.RedisClient
}

func (i appInfra) Close() {
	if i.redis != nil {
		_ = i.redis.Close()
	}
	if i.db != nil {
		_ = i.db.Close()
	}
}

type appRepositories struct {
	users          *sysuser.SysUserRepository
	roles          *role.RoleRepository
	menus          *menu.MenuRepository
	tokenBlacklist *infra.TokenBlacklistRepository
}

type serviceDeps struct {
	Cfg      *config.Config
	Tokens   token.Manager
	Policies policy.Manager
}

type appServices struct {
	health       *health.Service
	users        *sysuser.Service
	roles        *role.Service
	menus        *menu.Service
	tokenManager token.Manager
}

func NewApp(cfg *config.Config) (*App, error) {
	err := log.Init(cfg.Log)
	if err != nil {
		return nil, err
	}
	log.Info().Str("mode", cfg.Server.Mode).Msg("application starting")

	appReady := false
	var infra_ appInfra
	defer func() {
		if !appReady {
			infra_.Close()
		}
	}()

	infra_, err = newInfra(cfg)
	if err != nil {
		return nil, err
	}

	repos := newRepositories(infra_.db)
	policies, err := infra.NewCasbinManager(infra_.db.DB, repos.roles, cfg.Casbin.ModelPath)
	if err != nil {
		log.Error().Err(err).Msg("init casbin policy manager failed")
		return nil, err
	}

	tokenManager := token.NewManager(infra.NewJWTIssuer(&cfg.JWT), repos.tokenBlacklist)
	services := newServices(serviceDeps{
		Cfg:      cfg,
		Tokens:   tokenManager,
		Policies: policies,
	}, infra_, repos)

	init := initializer.NewInitializer(initializer.Deps{
		Cfg:      cfg,
		Users:    repos.users,
		Roles:    repos.roles,
		Menus:    repos.menus,
		Policies: policies,
	})
	if err := init.Run(context.Background()); err != nil {
		log.Error().Err(err).Msg("bootstrap initialization failed")
		return nil, err
	}

	if err := services.menus.LoadPermissionMap(context.Background()); err != nil {
		log.Error().Err(err).Msg("load permission map failed")
		return nil, err
	}

	docRegistry := apidoc.NewRegistry()
	adminRouter := newAdminRouter(services, policies, docRegistry)

	// Publisher 在路由注册前创建，作为 Router 传入 NewServer。
	// spec 惰性构建——路由注册完成后显式 Build() 以快速失败。
	pub := apidoc.NewPublisher(cfg.APIDoc, docRegistry)

	httpServer := server.NewServer(server.Config{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
		Mode: cfg.Server.Mode,
	}, adminRouter, pub)

	// 路由注册后主动构建 spec，确保启动时即可发现错误。
	if cfg.Server.Mode != gin.ReleaseMode {
		if err := pub.Build(); err != nil {
			log.Error().Err(err).Msg("build API documentation failed")
			return nil, err
		}
		log.Info().Msg("API documentation configured")
		log.Info().
			Str("url", fmt.Sprintf("http://%s:%d%s", cfg.Server.Host, cfg.Server.Port, pub.UIPath())).
			Msg("API docs (Redoc) available")
	}

	log.Info().Str("host", cfg.Server.Host).Int("port", cfg.Server.Port).Msg("server configured")

	appReady = true
	app := &App{HTTPServer: httpServer, DB: infra_.db, Redis: infra_.redis}
	app.startBlacklistCleanup(repos.tokenBlacklist)

	return app, nil
}

func newInfra(cfg *config.Config) (appInfra, error) {
	db, err := infra.NewDatabase(&cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("connect database failed")
		return appInfra{}, err
	}

	log.Info().Str("driver", cfg.Database.Driver).Msg("database connected")

	redisClient, err := infra.NewRedis(&cfg.Redis)
	if err != nil {
		log.Error().Err(err).Msg("connect redis failed")
		_ = db.Close()
		return appInfra{}, err
	}
	log.Info().Str("mode", cfg.Redis.Mode).Strs("addrs", cfg.Redis.Addrs).Msg("redis connected")

	return appInfra{db: db, redis: redisClient}, nil
}

// 数据仓库初始化
func newRepositories(db *infra.Database) appRepositories {
	return appRepositories{
		users:          sysuser.NewRepository(db.DB),
		roles:          role.NewRepository(db.DB),
		menus:          menu.NewRepository(db.DB),
		tokenBlacklist: infra.NewTokenBlacklistRepository(db.DB),
	}
}

func newServices(deps serviceDeps, infra_ appInfra, repos appRepositories) appServices {
	base := common.NewBaseService(deps.Cfg)
	menuSvc := menu.NewService(menu.Deps{Base: base, Repo: repos.menus})
	roleSvc := role.NewService(role.Deps{
		Base:   base,
		Repo:   repos.roles,
		Users:  repos.users,
		Menus:  menuSvc,
		Policy: deps.Policies,
	})

	return appServices{
		health: health.NewService(infra_.db, infra_.redis),
		users: sysuser.NewService(sysuser.Deps{
			Base:   base,
			Tokens: deps.Tokens,
			Users:  repos.users,
			Policy: deps.Policies,
			Roles:  roleSvc,
			Menus:  menuSvc,
		}),
		roles:        roleSvc,
		menus:        menuSvc,
		tokenManager: deps.Tokens,
	}
}

func newAdminRouter(services appServices, policies policy.Manager, docRegistry *apidoc.Registry) *router.AdminRouter {
	return router.NewAdminRouter(router.AdminRouterConfig{
		Health:      health.NewHandler(services.health),
		User:        sysuser.NewHandler(services.users),
		Role:        role.NewHandler(services.roles),
		Menu:        menu.NewHandler(services.menus),
		Tokens:      services.tokenManager,
		Policy:      policies,
		PermMap:     services.menus,
		DocRegistry: docRegistry,
	})
}

func (app *App) Cleanup() {
	if app.cleanupCancel != nil {
		app.cleanupCancel()
	}

	log.Info().Msg("application shutting down")

	if app.DB != nil {
		_ = app.DB.Close()
	}
	if app.Redis != nil {
		_ = app.Redis.Close()
	}
}

// 启动一个后台 goroutine，定期从数据库中移除过期的令牌黑名单
func (app *App) startBlacklistCleanup(repo *infra.TokenBlacklistRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	app.cleanupCancel = cancel

	go func() {
		ticker := time.NewTicker(blacklistCleanupInterval)
		defer ticker.Stop()

		// 在启动时立即运行一次
		if err := repo.DeleteExpired(ctx); err != nil {
			log.Error().Err(err).Msg("initial blacklist cleanup failed")
		} else {
			log.Debug().Msg("initial blacklist cleanup done")
		}

		for {
			select {
			case <-ticker.C:
				if err := repo.DeleteExpired(ctx); err != nil {
					log.Error().Err(err).Msg("periodic blacklist cleanup failed")
				} else {
					log.Debug().Msg("periodic blacklist cleanup done")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
