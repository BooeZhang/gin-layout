package model

import (
	"context"
	"golang.org/x/crypto/bcrypt"
)

// SysUser 系统用户
type SysUser struct {
	Model
	UserName  string `json:"username" gorm:"not null;unique;comment:用户名"`
	Password  string `json:"password" gorm:"not null;comment:密码"`
	Remark    string `json:"remark" gorm:"comment:备注"`
	HeaderImg string `json:"header_img" gorm:"comment:头像url"`
	IsSuper   bool   `json:"is_super" gorm:"not null;comment:是否是超级用户 0:不是 1:是"`
	IsActive  bool   `json:"is_active" gorm:"not null;comment:是否是激活状态 0:不是 1:是"`
}

type SysUserRepository interface {
	Create(c context.Context, user *SysUser) error
	Update(ctx context.Context, user *SysUser) error
	GetSysUserByName(c context.Context, name string) (*SysUser, error)
}

func (s *SysUser) TableName() string {
	return "sys_user"
}

// Encrypt encrypts the plain text with bcrypt.
func (s *SysUser) Encrypt() (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// Compare compares the encrypted text with the plain text if it's the same.
func (s *SysUser) Compare(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(password))
}
