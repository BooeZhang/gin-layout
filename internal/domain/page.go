package domain

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type PageRequest struct {
	Page     int
	PageSize int
}

func NewPageRequest(page, pageSize int) PageRequest {
	return PageRequest{Page: page, PageSize: pageSize}.Normalize()
}

func (p PageRequest) Normalize() PageRequest {
	if p.Page < DefaultPage {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 || p.PageSize > MaxPageSize {
		p.PageSize = DefaultPageSize
	}
	return p
}

func (p PageRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type PageResult[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

func NewPageResult[T any](items []T, total int64, page, pageSize int) PageResult[T] {
	if len(items) <= 0 {
		items = []T{}
	}
	return PageResult[T]{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

func (r PageResult[T]) TotalPages() int {
	if r.Total == 0 || r.PageSize == 0 {
		return 0
	}
	return int((r.Total + int64(r.PageSize) - 1) / int64(r.PageSize))
}

func (r PageResult[T]) HasMore() bool {
	return r.Page < r.TotalPages()
}
