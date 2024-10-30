package request

// ByID 根据 id 获取数据
type ByID struct {
	ID uint `json:"id" form:"id"` // 主键ID
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

// Safety 安全边界校验
func (p *PageInfo) Safety() {
	if p.Page <= 1 {
		p.Page = 0
	}

	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 50 {
		p.PageSize = 50
	}
}

// Offset 分页 offset 计算
func (p *PageInfo) Offset() int {
	if p.Page > 0 {
		return (p.Page - 1) * p.PageSize
	}
	return p.Page
}

// ListMeta 列表资源元数据
type ListMeta[T any] struct {
	PageInfo
	TotalCount int64 `json:"total,omitempty"`
	Items      []T   `json:"data"`
}

func NewListMeta[T any](data []T, pageInfo PageInfo, total int64) ListMeta[T] {
	if pageInfo.Page <= 1 {
		pageInfo.Page = 1
	}
	return ListMeta[T]{
		PageInfo:   pageInfo,
		Items:      data,
		TotalCount: total,
	}
}
