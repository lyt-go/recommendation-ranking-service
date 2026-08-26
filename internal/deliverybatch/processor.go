package deliverybatch

import (
	"recommendation/internal/deliveryrepo"
	"recommendation/internal/exportresource"
)

type Audit struct{ Entries []string }

func (a *Audit) Record(v string) { a.Entries = append(a.Entries, v) }

type Processor struct {
	Pool  *exportresource.Pool
	Repo  *deliveryrepo.Repo
	Audit *Audit
}

// Export 批量导出推荐结果。
//
// 语义约定：
//   - 成功：整批全部提交、释放导出句柄、审计记录 success。
//   - 写入失败：保留原始错误、整批回滚、审计记录 failure。
//
// 为此，导出句柄按“整批”获取（一次 Acquire，全程持有），而非每条记录
// 单独占用，从而避免在句柄池容量较小、仅因条目数超过容量时被误判为失败；
// Commit/Rollback 显式调用，确保返回的错误是真实的写入错误，而不是被
// 延迟提交静默吞掉。
func (p *Processor) Export(ids []string) (err error) {
	// 整批共享一个导出句柄：任意成功路径都需在退出前释放。
	lease, err := p.Pool.Acquire()
	if err != nil {
		// 句柄耗尽属于资源层面的失败，同样不应提前记成功。
		p.Audit.Record("failure")
		return err
	}
	defer lease.Release()

	tx := p.Repo.Begin()

	// 先写入全部条目，任一写入失败则保留原始错误、整批回滚、记失败。
	for _, id := range ids {
		if writeErr := tx.Write(id); writeErr != nil {
			_ = tx.Rollback()
			p.Audit.Record("failure")
			return writeErr
		}
	}

	// 全部写入成功后再提交；提交失败同样整批未生效，记失败。
	if commitErr := tx.Commit(); commitErr != nil {
		p.Audit.Record("failure")
		return commitErr
	}

	p.Audit.Record("success")
	return nil
}
