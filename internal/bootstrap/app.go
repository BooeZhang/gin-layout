package bootstrap

import (
	"context"
	"fmt"
	"time"

	"gin-layout/config"
	"gin-layout/internal/apidoc"
	"gin-layout/internal/auth"
	"gin-layout/internal/bootstrap/initializer"
	"gin-layout/internal/common"
	"gin-layout/internal/health"
	"gin-layout/internal/infra"
	"gin-layout/internal/menu"
	"gin-layout/internal/role"
	"gin-layout/internal/router"
	"gin-layout/internal/server"
	"gin-layout/internal/user"
)

type App struct {
	HTTPServer    *server.Server
	DB            *infra.Database
	Redis         *infra.RedisClient
	Logger        *infra.Logger
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
	users          *user.PGRepository
	roles          *role.PGRepository
	menus          *menu.PGRepository
	tokenBlacklist *infra.TokenBlacklistRepository
}

type appServices struct {
	health    *health.Service
	auth      *auth.Service
	users     *user.Service
	roles     *role.Service
	menus     *menu.Service
	tokens    common.TokenManager
	passwords common.PasswordHasher
}

func NewApp(cfg *config.Config) (*App, error) {
	logger := infra.NewLogger(&cfg.Log)
	logger.Info().Str("mode", cfg.Server.Mode).Msg("application starting")

	appReady := false
	var infra_ appInfra
	defer func() {
		if !appReady {
			infra_.Close()
		}
	}()

	infra_, err := newInfra(cfg, logger)
	if err != nil {
		return nil, err
	}

	repos := newRepositories(infra_.db)
	policies, err := infra.NewCasbinManager(infra_.db.DB, repos.roles, cfg.Casbin.ModelPath, logger)
	if err != nil {
		logger.Error().Err(err).Msg("init casbin policy manager failed")
		return nil, err
	}

	services := newServices(cfg, logger, infra_, repos, policies)

	init := initializer.NewInitializer(cfg, repos.users, repos.roles, repos.menus, services.passwords, policies, logger)
	if err := init.Run(context.Background()); err != nil {
		logger.Error().Err(err).Msg("bootstrap initialization failed")
		return nil, err
	}

	if err := services.menus.LoadPermissionMap(context.Background()); err != nil {
		logger.Error().Err(err).Msg("load permission map failed")
		return nil, err
	}

	docRegistry := apidoc.NewRegistry()
	adminRouter := newAdminRouter(logger, services, policies, docRegistry)

	// Publisher 在路由注册前创建，作为 Router 传入 NewServer。
	// spec 惰性构建——路由注册完成后显式 Build() 以快速失败。
	pub := apidoc.NewPublisher(cfg.APIDoc, docRegistry)

	httpServer := server.NewServer(server.Config{
		Host: cfg.Server.Host,
		Port: cfg.Server.Port,
		Mode: cfg.Server.Mode,
	}, logger, adminRouter, pub)

	// 路由注册后主动构建 spec，确保启动时即可发现错误。
	if err := pub.Build(); err != nil {
		logger.Error().Err(err).Msg("build API documentation failed")
		return nil, err
	}
	logger.Info().Msg("API documentation configured")

	if cfg.APIDoc.Enabled {
		logger.Info().
			Str("url", fmt.Sprintf("http://%s:%d%s/index.html", cfg.Server.Host, cfg.Server.Port, pub.UIPath())).
			Msg("swagger UI available")
	}

	logger.Info().Str("host", cfg.Server.Host).Int("port", cfg.Server.Port).Msg("server configured")

	appReady = true
	app := &App{HTTPServer: httpServer, DB: infra_.db, Redis: infra_.redis, Logger: logger}
	app.startBlacklistCleanup(repos.tokenBlacklist)
	return app, nil
}

func newInfra(cfg *config.Config, logger *infra.Logger) (appInfra, error) {
	db, err := infra.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Error().Err(err).Msg("connect database failed")
		return appInfra{}, err
	}
	logger.Info().Str("driver", cfg.Database.Driver).Msg("database connected")

	if err := db.Migrate(
		&user.UserModel{},
		&user.UserRoleModel{},
		&role.RoleModel{},
		&role.RoleMenuModel{},
		&menu.MenuModel{},
		&infra.TokenBlacklistModel{},
	); err != nil {
		logger.Error().Err(err).Msg("database migration failed")
		_ = db.Close()
		return appInfra{}, err
	}
	logger.Info().Msg("database migration completed")

	redisClient, err := infra.NewRedis(&cfg.Redis)
	if err != nil {
		logger.Error().Err(err).Msg("connect redis failed")
		_ = db.Close()
		return appInfra{}, err
	}
	logger.Info().Str("mode", cfg.Redis.Mode).Strs("addrs", cfg.Redis.Addrs).Msg("redis connected")

	return appInfra{db: db, redis: redisClient}, nil
}

func newRepositories(db *infra.Database) appRepositories {
	return appRepositories{
		users:          user.NewRepository(db.DB),
		roles:          role.NewRepository(db.DB),
		menus:          menu.NewRepository(db.DB),
		tokenBlacklist: infra.NewTokenBlacklistRepository(db.DB),
	}
}

func newServices(
	cfg *config.Config,
	logger *infra.Logger,
	infra_ appInfra,
	repos appRepositories,
	policies common.PolicyManager,
) appServices {
	passwords := infra.NewBcryptHasher()
	tokens := infra.NewJWTIssuer(&cfg.JWT)
	baseSvc := common.NewBaseService(cfg, logger)
	tokenManager := common.NewTokenManager(tokens, repos.tokenBlacklist)
	menuSvc := menu.NewService(baseSvc, repos.menus)
	roleSvc := role.NewService(baseSvc, repos.roles, repos.users, menuSvc, policies)

	return appServices{
		health:    health.NewService(infra_.db, infra_.redis),
		auth:      auth.NewService(baseSvc, repos.users, passwords, tokenManager),
		users:     user.NewService(baseSvc, repos.users, policies, passwords, roleSvc, menuSvc),
		roles:     roleSvc,
		menus:     menuSvc,
		tokens:    tokenManager,
		passwords: passwords,
	}
}

func newAdminRouter(logger *infra.Logger, services appServices, policies common.PolicyManager, docRegistry *apidoc.Registry) *router.AdminRouter {
	return router.NewAdminRouter(router.AdminRouterConfig{
		Auth:        auth.NewHandler(services.auth),
		Health:      health.NewHandler(services.health),
		User:        user.NewHandler(services.users),
		Role:        role.NewHandler(services.roles),
		Menu:        menu.NewHandler(services.menus),
		Tokens:      services.tokens,
		Policy:      policies,
		PermMap:     services.menus,
		Log:         logger,
		DocRegistry: docRegistry,
	})
}

func (app *App) Cleanup() {
	if app.cleanupCancel != nil {
		app.cleanupCancel()
	}
	if app.Logger != nil {
		app.Logger.Info().Msg("application shutting down")
	}
	if app.DB != nil {
		_ = app.DB.Close()
	}
	if app.Redis != nil {
		_ = app.Redis.Close()
	}
}

const blacklistCleanupInterval = 1 * time.Hour

// 启动一个后台 goroutine，定期从数据库中移除过期的令牌黑名单
func (app *App) startBlacklistCleanup(repo *infra.TokenBlacklistRepository) {
	ctx, cancel := context.WithCancel(context.Background())
	app.cleanupCancel = cancel

	go func() {
		ticker := time.NewTicker(blacklistCleanupInterval)
		defer ticker.Stop()

		// 在启动时立即运行一次
		if err := repo.DeleteExpired(ctx); err != nil {
			app.Logger.Error().Err(err).Msg("initial blacklist cleanup failed")
		} else {
			app.Logger.Debug().Msg("initial blacklist cleanup done")
		}

		for {
			select {
			case <-ticker.C:
				if err := repo.DeleteExpired(ctx); err != nil {
					app.Logger.Error().Err(err).Msg("periodic blacklist cleanup failed")
				} else {
					app.Logger.Debug().Msg("periodic blacklist cleanup done")
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
