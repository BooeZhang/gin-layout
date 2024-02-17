package config

// MongoConf mongo配置项
type MongoConf struct {
	Uri        string `json:"uri" mapstructure:"uri"`
	ReplicaSet string `yaml:"replicaset" mapstructure:"replicaset"`
	PoolLimit  int    `yaml:"poolLimit"  mapstructure:"poolLimit"`
}
