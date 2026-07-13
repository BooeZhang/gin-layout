package config

// Config 配置
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	Casbin      CasbinConfig      `mapstructure:"casbin"`
	Log         LogConfig         `mapstructure:"log"`
	Initializer InitializerConfig `mapstructure:"initializer"`
	APIDoc      APIDocConfig      `mapstructure:"apidoc"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`

	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`

	SSLMode  string `mapstructure:"sslmode"`
	Timezone string `mapstructure:"timezone"`
	DBFile   string `mapstructure:"dbFile"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addrs    []string `mapstructure:"addrs"`
	Mode     string   `mapstructure:"mode"`
	Password string   `mapstructure:"password"`
	DB       int      `mapstructure:"db"`
	PoolSize int      `mapstructure:"poolSize"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret         string `mapstructure:"secret"`
	AccessExpired  int    `mapstructure:"accessExpired"`
	RefreshExpired int    `mapstructure:"refreshExpired"`
}

// CasbinConfig Casbin配置
type CasbinConfig struct {
	ModelPath string `mapstructure:"modelPath"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"outputPath"`
}

// InitializerConfig 启动初始化配置
type InitializerConfig struct {
	AdminAccount  string `mapstructure:"adminAccount"`
	AdminPassword string `mapstructure:"adminPassword"`
	AdminRoleName string `mapstructure:"adminRoleName"`
	AdminRoleCode string `mapstructure:"adminRoleCode"`
	MenuJSON      string `mapstructure:"MenuJSON"`
}

// APIDocConfig API 文档配置
type APIDocConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	Strict      bool   `mapstructure:"strict"`
	JSONPath    string `mapstructure:"jsonPath"`
	UIPath      string `mapstructure:"uiPath"`
	Title       string `mapstructure:"title"`
	Version     string `mapstructure:"version"`
	Description string `mapstructure:"description"`
	Host        string `mapstructure:"host"`
	BasePath    string `mapstructure:"basePath"`
}
