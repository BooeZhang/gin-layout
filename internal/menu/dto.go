package menu

import "gin-layout/internal/domain"

type CreateMenuReq struct {
	ParentID   *int64          `json:"parent_id"`
	Name       string          `json:"name" binding:"required"`
	Code       *string         `json:"code" binding:"required"`
	Type       domain.MenuType `json:"type" binding:"required"`
	Path       string          `json:"path" binding:"required"`
	Redirect   string          `json:"redirect"`
	Component  string          `json:"component" binding:"required"`
	Icon       string          `json:"icon"`
	ActiveMenu string          `json:"active_menu"`
	Link       string          `json:"link"`
	Query      string          `json:"query"`
	Remark     string          `json:"remark"`
	Sort       int             `json:"sort"`
	Level      int             `json:"level"`
	Hidden     *bool           `json:"hidden"`
	Cache      *bool           `json:"cache"`
	Affix      *bool           `json:"affix"`
	Breadcrumb *bool           `json:"breadcrumb"`
	AlwaysShow *bool           `json:"always_show"`
	External   *bool           `json:"external"`
	Iframe     *bool           `json:"iframe"`
	Enabled    *bool           `json:"enabled"`
	Method     string          `json:"method"`
	APIPath    string          `json:"apiPath"`
	PermCode   *string         `json:"permCode"`
}

type CreateMenuRes struct {
	ID int64 `json:"id"`
}

type UpdateMenuReq struct {
	MenuID     int64            `json:"-"`
	ParentID   *int64           `json:"parent_id"`
	Name       *string          `json:"name"`
	RouteName  *string          `json:"route_name"`
	Code       *string          `json:"code"`
	Type       *domain.MenuType `json:"type"`
	Path       *string          `json:"path"`
	Method     *string          `json:"method"`
	Redirect   *string          `json:"redirect"`
	Component  *string          `json:"component"`
	Icon       *string          `json:"icon"`
	ActiveMenu *string          `json:"active_menu"`
	Link       *string          `json:"link"`
	Query      *string          `json:"query"`
	Remark     *string          `json:"remark"`
	Sort       *int             `json:"sort"`
	Level      *int             `json:"level"`
	Hidden     *bool            `json:"hidden"`
	Cache      *bool            `json:"cache"`
	Affix      *bool            `json:"affix"`
	Breadcrumb *bool            `json:"breadcrumb"`
	AlwaysShow *bool            `json:"always_show"`
	External   *bool            `json:"external"`
	Iframe     *bool            `json:"iframe"`
	Enabled    *bool            `json:"enabled"`
	APIPath    *string          `json:"apiPath"`
	PermCode   *string          `json:"permCode"`
}

type UpdateMenuRes struct{}
