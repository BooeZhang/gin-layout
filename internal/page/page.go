package page

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Request struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func NewRequest(pageNum, pageSize int) Request {
	return Request{Page: pageNum, PageSize: pageSize}.Normalize()
}

func (p Request) Normalize() Request {
	if p.Page < DefaultPage {
		p.Page = DefaultPage
	}
	if p.PageSize < 1 || p.PageSize > MaxPageSize {
		p.PageSize = DefaultPageSize
	}
	return p
}

func (p Request) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type Result[T any] struct {
	Items    []T   `json:"items"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

func NewResult[T any](items []T, total int64, pageNum, pageSize int) Result[T] {
	if len(items) <= 0 {
		items = []T{}
	}
	return Result[T]{
		Items:    items,
		Total:    total,
		Page:     pageNum,
		PageSize: pageSize,
	}
}

func (r Result[T]) TotalPages() int {
	if r.Total == 0 || r.PageSize == 0 {
		return 0
	}
	return int((r.Total + int64(r.PageSize) - 1) / int64(r.PageSize))
}

func (r Result[T]) HasMore() bool {
	return r.Page < r.TotalPages()
}
