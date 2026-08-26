package provideradapter

import "fmt"

type RejectedError struct{ Message string }

func (e RejectedError) Error() string { return e.Message }

type TemporaryError struct{ Message string }

func (e TemporaryError) Error() string   { return e.Message }
func (e TemporaryError) Temporary() bool { return true }

type Backend interface{ Fetch() (string, error) }
type Adapter struct{ backend Backend }

func New(backend Backend) *Adapter { return &Adapter{backend: backend} }
func (a *Adapter) Fetch() (string, error) {
	value, err := a.backend.Fetch()
	if err != nil {
		return "", fmt.Errorf("provider failed: %v", err)
	}
	return value, nil
}
