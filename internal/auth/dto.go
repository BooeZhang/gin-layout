package auth

type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type (
	LogoutReq struct{}
	LogoutRes struct{}
)

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenRes struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
