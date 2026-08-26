// Package rankdispatch 把"异步排序"的请求派发到后台 goroutine。
//
// 它依赖 requestscope.Pool 复用请求身份对象。正因为对象被复用，派发逻辑
// 必须保证：每条请求在跨过 goroutine 边界后仍持有自己的身份，不会因为
// 相邻请求复用了同一个 *Scope 而串号。具体做法是——在把 *Scope 归还池子
// 之前，先取出与池化对象完全解耦的快照值，异步任务只持有快照、绝不触碰
// 池化的 *Scope 本身。
package rankdispatch

import "recommendation/internal/requestscope"

// Dispatcher 负责把排序请求异步派发。
//
// gate 是一个一次性门控信道：收到信号后才放行后台排序回调。在门控打开
// 之前，请求身份必须已经固化为独立快照，避免与后续请求复用同一 *Scope
// 时发生串读。
type Dispatcher struct {
	pool *requestscope.Pool
	gate <-chan struct{}
}

// New 创建一个派发器。
func New(pool *requestscope.Pool, gate <-chan struct{}) *Dispatcher {
	return &Dispatcher{pool: pool, gate: gate}
}

// Dispatch 派发一次异步排序。
//
// 调用顺序刻意安排为：
//  1. Acquire 取出（或复用）一个 *Scope，写入本次请求身份；
//  2. 立刻 Snapshot，把身份固化为与池化对象解耦的值拷贝；
//  3. Release 清空并归还 *Scope；
//  4. 启动 goroutine，gate 打开后回调只持有第 2 步的快照。
//
// 这样即便 req-b 紧随 req-a 到来、复用了同一个 *Scope 并写入 user-b，
// req-a 的回调拿到的也始终是它自己的快照 user-a，不会串号。
func (d *Dispatcher) Dispatch(requestID, userID string, done func(requestscope.Snapshot)) {
	s := d.pool.Acquire(requestID, userID, nil)
	snap := s.Snapshot()
	d.pool.Release(s)
	go func() {
		<-d.gate
		done(snap)
	}()
}
