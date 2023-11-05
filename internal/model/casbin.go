package model

// UserRule 用户权限
type UserRule struct {
	Model
	Ptype string `json:"ptype"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
}

func (u *UserRule) TableName() string {
	return "casbin_rule"
}
