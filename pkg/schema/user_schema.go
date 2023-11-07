package schema

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginRes struct {
	Token string `json:"token"`
}
