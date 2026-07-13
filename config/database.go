package config

import "fmt"

var dsnBuilders = map[string]func(DatabaseConfig) string{
	"mysql":     buildMySQLDSN,
	"postgres":  buildPostgresDSN,
	"sqlite":    buildSQLiteDSN,
	"sqlserver": buildSQLServerDSN,
}

// BuildDatabaseDSN 构建数据库DSN
func BuildDatabaseDSN(cfg DatabaseConfig) (string, error) {
	builder, ok := dsnBuilders[cfg.Driver]
	if !ok {
		return "", fmt.Errorf("unsupported database driver: %q", cfg.Driver)
	}
	return builder(cfg), nil
}

func buildMySQLDSN(cfg DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func buildPostgresDSN(cfg DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func buildSQLiteDSN(cfg DatabaseConfig) string {
	return cfg.DBFile
}

func buildSQLServerDSN(cfg DatabaseConfig) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
