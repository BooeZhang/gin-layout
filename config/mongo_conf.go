package config

// MongoConf mongo配置项
type MongoConf struct {
	Host       string `json:"host,omitempty"                     mapstructure:"host"`
	Username   string `json:"username,omitempty"                 mapstructure:"username"`
	Password   string `json:"password"                           mapstructure:"password"`
	Database   string `json:"database"                           mapstructure:"database"`
	ReplicaSet string `yaml:"replicaset" mapstructure:"replicaset"`
	PoolLimit  int    `yaml:"poolLimit"  mapstructure:"poolLimit"`
}

// NewMongoConf 创建一个默认的mongo配置项.
func NewMongoConf() *MongoConf {
	return &MongoConf{
		Host:       "127.0.0.1:27017",
		Username:   "",
		Password:   "",
		Database:   "fighting_boss",
		ReplicaSet: "",
		PoolLimit:  500,
	}
}
