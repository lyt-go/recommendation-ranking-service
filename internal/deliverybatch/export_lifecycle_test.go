package deliverybatch_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"recommendation/internal/deliverybatch"
	"recommendation/internal/deliveryrepo"
	"recommendation/internal/exportresource"
)

func TestBatchExportReleasesHandlesAndKeepsTransactionOutcome(t *testing.T) {
	pool := exportresource.NewPool(2)
	repo := &deliveryrepo.Repo{}
	audit := &deliverybatch.Audit{}
	p := &deliverybatch.Processor{Pool: pool, Repo: repo, Audit: audit}
	if err := p.Export([]string{"rec-a", "rec-b", "rec-c"}); err != nil {
		t.Fatalf("three-record export returned %v", err)
	}
	if !reflect.DeepEqual(repo.Committed, []string{"rec-a", "rec-b", "rec-c"}) {
		t.Errorf("committed recommendations = %v, want all three", repo.Committed)
	}
	if pool.Active() != 0 {
		t.Errorf("active export handles = %d, want 0", pool.Active())
	}
	if !reflect.DeepEqual(audit.Entries, []string{"success"}) {
		t.Errorf("successful export audit = %v", audit.Entries)
	}

	pool2 := exportresource.NewPool(2)
	repo2 := &deliveryrepo.Repo{FailOn: "rec-b"}
	audit2 := &deliverybatch.Audit{}
	p2 := &deliverybatch.Processor{Pool: pool2, Repo: repo2, Audit: audit2}
	if err := p2.Export([]string{"rec-a", "rec-b"}); !errors.Is(err, deliveryrepo.ErrWrite) {
		t.Fatalf("failed export error = %v, want delivery write failure", err)
	}
	if len(repo2.Committed) != 0 || len(audit2.Entries) != 1 || !strings.HasPrefix(audit2.Entries[0], "fail") {
		t.Fatalf("failed export committed=%v audit=%v, want rollback and one failed audit", repo2.Committed, audit2.Entries)
	}
}
