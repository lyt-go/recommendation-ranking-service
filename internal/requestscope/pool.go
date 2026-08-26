package requestscope

import "sync"

// Scope 是单个请求在其整个处理生命周期内携带的请求身份与上下文。
//
// Scope 由 Pool 维护，是可复用的池化对象：一个 *Scope 在归还后会再次被
// Acquire 取出并复用。因此使用方必须遵守两条不变量：
//   - 调用 Release 之后，绝不能再读写该 *Scope；
//   - 异步任务若需要请求身份，必须先调用 Snapshot 取出独立的快照值，
//     再把 *Scope 归还；异步任务只持有快照，不触碰池化的 *Scope 本身。
//
// 这两点共同保证：即便相邻两次请求复用了同一个 *Scope，各自的请求身份
// （RequestID / UserID / Tags）也不会发生串读。
type Scope struct {
	RequestID string
	UserID    string
	Tags      []string
}

// Snapshot 是 Scope 的值拷贝，与池化的 *Scope 完全解耦。
// Tags 会被整体复制，因此在快照生成之后对原 *Scope 的任何复用、改写
// 都不会影响已取出的快照。
type Snapshot struct {
	RequestID string
	UserID    string
	Tags      []string
}

// Pool 用 sync.Pool 复用 *Scope，减少每请求分配。
type Pool struct{ pool sync.Pool }

// NewPool 创建一个 Scope 池。
func NewPool() *Pool {
	p := &Pool{}
	p.pool.New = func() any { return &Scope{} }
	return p
}

// Acquire 取出一个 *Scope 并写入本次请求的身份。
//
// tags 永远在一张全新的底层数组上构建：既不会与上一任请求的 Tags 共享
// 底层数组，也不会在复用时累积残留的 tag。
func (p *Pool) Acquire(requestID, userID string, tags []string) *Scope {
	s := p.pool.Get().(*Scope)
	s.RequestID = requestID
	s.UserID = userID
	s.Tags = append([]string(nil), tags...)
	return s
}

// Release 清空所有身份字段后将 *Scope 归还。
//
// 归还时强制清零 RequestID / UserID / Tags：即使后续某次 Acquire 传入了
// 空值，复用对象也不会泄漏上一次请求的身份。
func (p *Pool) Release(s *Scope) {
	if s == nil {
		return
	}
	s.RequestID = ""
	s.UserID = ""
	s.Tags = nil
	p.pool.Put(s)
}

// Snapshot 返回与池化 *Scope 完全解耦的值拷贝。
//
// Tags 做整片 copy，使得快照在生成之后不再受 *Scope 后续复用的影响，
// 异步任务可以安全地长期持有快照。
func (s *Scope) Snapshot() Snapshot {
	tags := make([]string, len(s.Tags))
	copy(tags, s.Tags)
	return Snapshot{RequestID: s.RequestID, UserID: s.UserID, Tags: tags}
}
