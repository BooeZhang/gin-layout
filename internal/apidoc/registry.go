package apidoc

import "sync"

// Registry 是在路由注册期间收集端点元数据的线程安全存储。
// 路由注册在启动时（单线程）完成，因此基于指针的修改是安全的
// —— setter 在初始推送后原地修改记录。
type Registry struct {
	mu    sync.RWMutex
	items []*EndpointRecord
}

// NewRegistry 创建一个空的 Registry。
func NewRegistry() *Registry {
	return &Registry{}
}

// Add 记录一个端点。存储指针以便后续通过 Route setter 对记录的修改对 builder 可见。线程安全。
func (r *Registry) Add(ep *EndpointRecord) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = append(r.items, ep)
}

// Items 返回所有已记录端点的快照。线程安全。
func (r *Registry) Items() []*EndpointRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*EndpointRecord, len(r.items))
	copy(out, r.items)
	return out
}
