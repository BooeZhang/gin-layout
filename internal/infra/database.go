package infra

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gin-layout/config"
)

var ErrDbIsNil = fmt.Errorf("database is nil")

type Database struct {
	DB *gorm.DB
}

var dialectorOpeners = map[string]func(string) gorm.Dialector{
	"mysql":     mysql.Open,
	"postgres":  postgres.Open,
	"sqlite":    sqlite.Open,
	"sqlserver": sqlserver.Open,
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	opener, ok := dialectorOpeners[cfg.Driver]
	if !ok {
		return nil, fmt.Errorf("unsupported database driver: %q", cfg.Driver)
	}
	dialector := opener(cfg.DSN)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("connect database ping failed: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Migrate(models ...any) error {
	return d.DB.AutoMigrate(models...)
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Ping(ctx context.Context) error {
	if d == nil || d.DB == nil {
		return ErrDbIsNil
	}
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
