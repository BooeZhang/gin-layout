package schema

// ById 根据 id 获取数据
type ById struct {
	ID uint32 `json:"id" form:"id"` // 主键ID
}

// ByIds 根据 id 列表获取数据
type ByIds struct {
	Ids []uint32 `json:"ids" form:"ids"`
}

// PageInfo 分页参数
type PageInfo struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}
