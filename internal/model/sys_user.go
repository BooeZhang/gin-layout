package model

import "golang.org/x/crypto/bcrypt"

// SysUser 系统用户
type SysUser struct {
	Model
	Account  string `json:"account" gorm:"not null;unique;comment:账户名"`
	Password string `json:"password" gorm:"not null;comment:密码"`
	Remark   string `json:"remark" gorm:"comment:备注"`
	IsActive bool   `json:"is_active" gorm:"not null;comment:是否是激活状态 0:不是 1:是"`
}

func (s *SysUser) TableName() string {
	return "sys_user"
}

// Encrypt 加密密码.
func (s *SysUser) Encrypt() (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare 密码比较
func (s *SysUser) Compare(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
}
