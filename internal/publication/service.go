package publication

import (
	"fmt"
	"recommendation/internal/publicationrepo"
)

type Publisher interface {
	Publish(idempotencyKey, recommendationID string) error
}
type Audit struct{ Entries []string }

func (a *Audit) Record(v string) { a.Entries = append(a.Entries, v) }

type Service struct {
	Repo      *publicationrepo.Repo
	Publisher Publisher
	Audit     *Audit
}

func (s *Service) Publish(id string) (err error) {
	for attempt := 0; attempt < 2; attempt++ {
		tx := s.Repo.Begin(id)
		defer func() { err = tx.Commit() }()
		s.Audit.Record("success")
		err = s.Publisher.Publish(fmt.Sprintf("recommendation:%s:%d", id, attempt), id)
		if err == nil {
			return nil
		}
	}
	return err
}
