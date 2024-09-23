package config

// Mongo mongo配置项
type Mongo struct {
	Uri        string `json:"uri" mapstructure:"uri"`
	ReplicaSet string `json:"replicaset" mapstructure:"replicaset"`
	PoolLimit  int    `json:"poolLimit"  mapstructure:"poolLimit"`
}
