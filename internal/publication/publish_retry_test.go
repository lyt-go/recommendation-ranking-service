package publication_test

import (
	"errors"
	"recommendation/internal/publication"
	"recommendation/internal/publicationrepo"
	"reflect"
	"testing"
)

type temporaryError struct{}

func (temporaryError) Error() string   { return "publisher temporarily unavailable" }
func (temporaryError) Temporary() bool { return true }

type publisher struct {
	calls  int
	unique map[string]bool
}

func (p *publisher) Publish(key, id string) error {
	p.calls++
	p.unique[key] = true
	if p.calls == 1 {
		return temporaryError{}
	}
	return nil
}

func TestTemporaryPublishFailureRetriesWithoutDuplicateCommitOrEvent(t *testing.T) {
	repo := &publicationrepo.Repo{}
	pub := &publisher{unique: make(map[string]bool)}
	audit := &publication.Audit{}
	svc := &publication.Service{Repo: repo, Publisher: pub, Audit: audit}
	if err := svc.Publish("rec-42"); err != nil {
		t.Fatalf("publication returned %v", err)
	}
	if pub.calls != 2 {
		t.Errorf("publisher calls = %d, want 2", pub.calls)
	}
	if len(pub.unique) != 1 {
		t.Errorf("published event identities = %d, want 1 stable event", len(pub.unique))
	}
	if repo.Commits != 1 || !reflect.DeepEqual(repo.IDs, []string{"rec-42"}) {
		t.Errorf("repository commits=%d ids=%v, want one commit", repo.Commits, repo.IDs)
	}
	if !reflect.DeepEqual(audit.Entries, []string{"success"}) {
		t.Errorf("audit entries=%v, want one final success", audit.Entries)
	}
	if errors.Is(nil, temporaryError{}) {
		t.Fatal("unreachable")
	}
}
