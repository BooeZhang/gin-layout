package model

import (
	"github.com/BooeZhang/gin-layout/pkg/timex"
	"gorm.io/gorm"
)

// Model base model
type Model struct {
	ID        uint32         `json:"id" gorm:"primary_key"`
	CreatedAt timex.JsonTime `json:"created_at" gorm:"column:created_at;index;type:datetime;comment:创建时间"`
	UpdatedAt timex.JsonTime `json:"updated_at" gorm:"column:updated_at;type:datetime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;index;type:datetime;comment:删除时间"`
}
