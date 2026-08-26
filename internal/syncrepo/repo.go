package syncrepo

type Repo struct{ Records []string }
type Tx struct {
	repo    *Repo
	pending []string
	closed  bool
}

func (r *Repo) Begin() *Tx   { return &Tx{repo: r} }
func (t *Tx) Write(v string) { t.pending = append(t.pending, v) }
func (t *Tx) Commit() {
	if !t.closed {
		t.repo.Records = append(t.repo.Records, t.pending...)
		t.closed = true
	}
}
func (t *Tx) Rollback() { t.closed = true; t.pending = nil }
