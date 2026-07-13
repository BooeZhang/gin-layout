package domain

import "time"

type User struct {
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
