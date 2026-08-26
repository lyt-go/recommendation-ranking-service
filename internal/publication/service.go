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
	// 同一次发布的所有重试共享同一幂等键，代表一个稳定事件身份。
	key := fmt.Sprintf("recommendation:%s", id)
	for attempt := 0; attempt < 2; attempt++ {
		tx := s.Repo.Begin(id)
		// 只有发布端真正成功后才提交推荐记录并记审计，失败则回滚重试。
		if err = s.Publisher.Publish(key, id); err == nil {
			s.Audit.Record("success")
			return tx.Commit()
		}
		_ = tx.Rollback()
	}
	return err
}
