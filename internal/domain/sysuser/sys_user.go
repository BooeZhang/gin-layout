package sysuser

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SysUser 系统用户
type SysUser struct {
	ID           int64      `gorm:"primaryKey"`
	Account      string     `gorm:"uniqueIndex;size:50;not null"`
	PasswordHash string     `gorm:"column:password;size:255;not null"`
	NickName     string     `gorm:"size:100"`
	Email        string     `gorm:"size:100;index"`
	Phone        string     `gorm:"size:20;index"`
	Avatar       string     `gorm:"size:500"`
	Enabled      bool       `gorm:"default:true"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	RoleIDs      []int64    `gorm:"-"`
}

func (*SysUser) TableName() string { return "sys_user" }

// SysUserRole 用户-角色关联表。
type SysUserRole struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey;index"`
}

func (SysUserRole) TableName() string { return "sys_user_roles" }

func (u *SysUser) PwdHash() error {
	if u.PasswordHash == "" {
		return errors.New("password hash is empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}

func (u *SysUser) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}
