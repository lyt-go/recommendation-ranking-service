package deliverytask

const (
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
)

type Task struct {
	ID      string
	Version int
	Status  string
}
