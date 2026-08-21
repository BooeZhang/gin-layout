package menu

import "time"

type Type string

const (
	TypeCatalog Type = "catalog"
	TypeMenu    Type = "menu"
	TypeButton  Type = "button"
)

// Menu 菜单,同时作为领域对象与 GORM 持久化模型。
type Menu struct {
	ID         int64     `gorm:"primaryKey"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	ParentID   *int64    `gorm:"index"`
	Name       string    `gorm:"size:100;index;not null"`
	Code       *string   `gorm:"size:100;uniqueIndex:uk_menu_code"`
	MenuType   Type      `gorm:"size:20;default:'menu';index"`
	Path       string    `gorm:"size:255;index"`
	Redirect   string    `gorm:"size:255"`
	Component  string    `gorm:"size:255"`
	Icon       string    `gorm:"size:100"`
	ActiveMenu string    `gorm:"column:active_menu;size:255"`
	Link       string    `gorm:"size:500"`
	Query      string    `gorm:"type:text"`
	Remark     string    `gorm:"size:255"`
	Sort       int       `gorm:"default:0;index"`
	Level      int       `gorm:"default:0"`
	Hidden     bool      `gorm:"default:false"`
	Cache      bool      `gorm:"default:true"`
	Affix      bool      `gorm:"default:false"`
	Breadcrumb bool      `gorm:"default:true"`
	AlwaysShow bool      `gorm:"column:always_show;default:false"`
	External   bool      `gorm:"default:false"`
	Iframe     bool      `gorm:"default:false"`
	Enabled    bool      `gorm:"default:true;index"`
	Method     string    `gorm:"size:10"`
	APIPath    string    `gorm:"size:255"`
	PermCode   *string   `gorm:"size:255;uniqueIndex:uk_menu_perm_code"`
}

func (Menu) TableName() string { return "menus" }
