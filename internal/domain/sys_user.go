package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SysUser 系统用户域
type SysUser struct {
	ID           int64
	Account      string
	PasswordHash string
	NickName     string
	Email        string
	Phone        string
	Avatar       string
	Enabled      bool
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	RoleIDs      []int64
}

// PwdHash 密码 hash
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

// ComparePassword 比较密码
func (u *SysUser) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}
