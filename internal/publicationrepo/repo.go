package publicationrepo

type Repo struct {
	Commits int
	IDs     []string
}
type Tx struct {
	repo   *Repo
	id     string
	closed bool
}

func (r *Repo) Begin(id string) *Tx { return &Tx{repo: r, id: id} }
func (t *Tx) Commit() error {
	if !t.closed {
		t.repo.Commits++
		t.repo.IDs = append(t.repo.IDs, t.id)
		t.closed = true
	}
	return nil
}
func (t *Tx) Rollback() error { t.closed = true; return nil }
