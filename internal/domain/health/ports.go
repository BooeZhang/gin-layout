package health

import "context"

// Pinger 可以报告其健康状态的组件的接口
type Pinger interface {
	Ping(ctx context.Context) error
}
