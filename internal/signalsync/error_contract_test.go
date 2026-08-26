package signalsync_test

import (
	"errors"
	"recommendation/internal/provideradapter"
	"recommendation/internal/signalsync"
	"recommendation/internal/syncrepo"
	"reflect"
	"testing"
)

type backend struct {
	errors []error
	value  string
	calls  int
}

func (b *backend) Fetch() (string, error) {
	b.calls++
	if len(b.errors) > 0 {
		err := b.errors[0]
		b.errors = b.errors[1:]
		if err != nil {
			return "", err
		}
	}
	return b.value, nil
}
func TestSignalSyncPreservesErrorClassAndCommitsOnlySuccessfulRetry(t *testing.T) {
	rejected := &backend{errors: []error{provideradapter.RejectedError{Message: "signal rejected"}, provideradapter.RejectedError{Message: "signal rejected"}}}
	rejectedRepo := &syncrepo.Repo{}
	rejectedSvc := signalsync.New(provideradapter.New(rejected), rejectedRepo)
	var rejection provideradapter.RejectedError
	if err := rejectedSvc.Sync(); !errors.As(err, &rejection) {
		t.Errorf("rejection error=%v, want typed rejection", err)
	}
	if rejected.calls != 1 || len(rejectedRepo.Records) != 0 {
		t.Errorf("rejected sync calls=%d records=%v, want one call and no writes", rejected.calls, rejectedRepo.Records)
	}
	temporary := &backend{errors: []error{provideradapter.TemporaryError{Message: "try later"}, nil}, value: "signal-7"}
	temporaryRepo := &syncrepo.Repo{}
	temporarySvc := signalsync.New(provideradapter.New(temporary), temporaryRepo)
	if err := temporarySvc.Sync(); err != nil {
		t.Fatalf("temporary sync returned %v", err)
	}
	if temporary.calls != 2 || !reflect.DeepEqual(temporaryRepo.Records, []string{"signal-7"}) {
		t.Fatalf("temporary sync calls=%d records=%v, want two calls and one committed signal", temporary.calls, temporaryRepo.Records)
	}
}
