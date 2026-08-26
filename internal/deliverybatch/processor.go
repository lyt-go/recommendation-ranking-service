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

func (p *Processor) Export(ids []string) (err error) {
	tx := p.Repo.Begin()
	defer func() { err = tx.Commit() }()
	p.Audit.Record("success")
	for _, id := range ids {
		lease, acquireErr := p.Pool.Acquire()
		if acquireErr != nil {
			return acquireErr
		}
		defer lease.Release()
		if writeErr := tx.Write(id); writeErr != nil {
			return writeErr
		}
	}
	return nil
}
