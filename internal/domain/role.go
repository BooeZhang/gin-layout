package domain

import "time"

type Role struct {
	ID          int64
	Name        string
	Code        string
	Description string
	Sort        int
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	MenuIDs     []int64
}
