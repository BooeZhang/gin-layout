package domain

type MenuType string

type MenuItem struct {
	ID         int64      `json:"id"`
	ParentID   *int64     `json:"parentId"`
	Name       string     `json:"name"`
	Code       string     `json:"code"`
	Type       string     `json:"type"`
	Path       string     `json:"path"`
	Redirect   string     `json:"redirect"`
	Component  string     `json:"component"`
	Icon       string     `json:"icon"`
	ActiveMenu string     `json:"activeMenu"`
	Link       string     `json:"link"`
	Query      string     `json:"query"`
	Remark     string     `json:"remark"`
	Sort       int        `json:"sort"`
	Level      int        `json:"level"`
	Hidden     bool       `json:"hidden"`
	Cache      bool       `json:"cache"`
	Affix      bool       `json:"affix"`
	Breadcrumb bool       `json:"breadcrumb"`
	AlwaysShow bool       `json:"alwaysShow"`
	External   bool       `json:"external"`
	Iframe     bool       `json:"iframe"`
	Enabled    bool       `json:"enabled"`
	Method     string     `json:"method"`
	APIPath    string     `json:"apiPath"`
	PermCode   string     `json:"permCode"`
	Children   []MenuItem `json:"children,omitempty"`
}

const (
	MenuTypeCatalog MenuType = "catalog"
	MenuTypeMenu    MenuType = "menu"
	MenuTypeButton  MenuType = "button"
)

type Menu struct {
	ID         int64
	ParentID   *int64
	Name       string
	Code       *string
	Type       MenuType
	Path       string
	Redirect   string
	Component  string
	Icon       string
	ActiveMenu string
	Link       string
	Query      string
	Remark     string
	Sort       int
	Level      int
	Hidden     bool
	Cache      bool
	Affix      bool
	Breadcrumb bool
	AlwaysShow bool
	External   bool
	Iframe     bool
	Enabled    bool
	Method     string
	APIPath    string
	PermCode   *string
	UpdatedAt  int64
	CreatedAt  int64
}
