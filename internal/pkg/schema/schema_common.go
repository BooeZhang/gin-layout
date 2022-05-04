package schema

type ById struct {
	ID uint32 `json:"id" form:"id"` // 主键ID
}

type ByIds struct {
	Ids []uint32 `json:"ids" form:"ids"`
}

type PageInfo struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}
