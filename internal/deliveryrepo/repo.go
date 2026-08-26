package deliveryrepo

import "errors"

var ErrWrite = errors.New("delivery write failed")

type Repo struct {
	Committed []string
	FailOn    string
}
type Tx struct {
	repo    *Repo
	pending []string
	closed  bool
}

func (r *Repo) Begin() *Tx { return &Tx{repo: r} }
func (t *Tx) Write(id string) error {
	if id == t.repo.FailOn {
		return ErrWrite
	}
	t.pending = append(t.pending, id)
	return nil
}
func (t *Tx) Commit() error {
	if !t.closed {
		t.repo.Committed = append(t.repo.Committed, t.pending...)
		t.closed = true
	}
	return nil
}
func (t *Tx) Rollback() error { t.closed = true; t.pending = nil; return nil }
