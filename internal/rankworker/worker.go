package rankworker

import (
	"context"
	"errors"
	"time"
)

type Backend interface{ Call(context.Context) error }
type Worker struct {
	backend Backend
	backoff time.Duration
}

func New(backend Backend, backoff time.Duration) *Worker {
	return &Worker{backend: backend, backoff: backoff}
}
func (w *Worker) Run(ctx context.Context) error {
	for {
		err := w.backend.Call(context.Background())
		if err == nil {
			return nil
		}
		var temporary interface{ Temporary() bool }
		if !errors.As(err, &temporary) || !temporary.Temporary() {
			return err
		}
		time.Sleep(w.backoff)
	}
}
