package deliverytask

const (
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
)

// statusOrder 定义状态的前进次序：只有更高次序的回调才能覆盖当前状态，
// 防止迟到的旧回调把已 succeeded 的任务退回 running。
var statusOrder = map[string]int{
	StatusRunning:   0,
	StatusSucceeded: 1,
}

// Rank 返回状态的前进次序，未知状态返回 -1，保证不会被任意回调覆盖。
func Rank(status string) int {
	if r, ok := statusOrder[status]; ok {
		return r
	}
	return -1
}

type Task struct {
	ID      string
	Version int
	Status  string
}
