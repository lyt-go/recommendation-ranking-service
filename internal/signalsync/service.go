package signalsync

import (
	"errors"
	"recommendation/internal/provideradapter"
	"recommendation/internal/syncrepo"
)

type Service struct {
	adapter *provideradapter.Adapter
	repo    *syncrepo.Repo
}

func New(adapter *provideradapter.Adapter, repo *syncrepo.Repo) *Service {
	return &Service{adapter: adapter, repo: repo}
}
func (s *Service) Sync() (err error) {
	for attempt := 0; attempt < 2; attempt++ {
		tx := s.repo.Begin()
		tx.Write("partial")
		defer tx.Commit()
		value, fetchErr := s.adapter.Fetch()
		if fetchErr != nil {
			err = fetchErr
			continue
		}
		tx.Write(value)
		return nil
	}
	return errors.New("signal retries exhausted: " + err.Error())
}
