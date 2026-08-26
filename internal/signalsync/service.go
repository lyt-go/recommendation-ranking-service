package signalsync

import (
	"errors"
	"recommendation/internal/provideradapter"
	"recommendation/internal/syncrepo"
)

const maxAttempts = 2

type Service struct {
	adapter *provideradapter.Adapter
	repo    *syncrepo.Repo
}

func New(adapter *provideradapter.Adapter, repo *syncrepo.Repo) *Service {
	return &Service{adapter: adapter, repo: repo}
}
func (s *Service) Sync() error {
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		value, fetchErr := s.adapter.Fetch()
		if fetchErr == nil {
			// 仅在成功时提交，且只写入真实信号，不写任何占位记录。
			tx := s.repo.Begin()
			tx.Write(value)
			tx.Commit()
			return nil
		}

		lastErr = fetchErr

		// 明确拒绝属于永久性错误：保留原始错误类型，立即返回，不重试、不写记录。
		var rejected provideradapter.RejectedError
		if errors.As(fetchErr, &rejected) {
			return fetchErr
		}

		// 非临时错误同样不应重试（如未知类型），原样返回，不写记录。
		if !isTemporary(fetchErr) {
			return fetchErr
		}
		// 临时失败：丢弃本次未提交的事务后重试。
	}
	return lastErr
}

// isTemporary 判断是否为可重试的临时错误。
func isTemporary(err error) bool {
	type temporary interface{ Temporary() bool }
	var t temporary
	if errors.As(err, &t) {
		return t.Temporary()
	}
	return false
}
