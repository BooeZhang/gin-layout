package config

// MongoConf mongo配置项
type MongoConf struct {
	Uri        string `json:"uri" mapstructure:"uri"`
	ReplicaSet string `json:"replicaset" mapstructure:"replicaset"`
	PoolLimit  int    `json:"poolLimit"  mapstructure:"poolLimit"`
}
