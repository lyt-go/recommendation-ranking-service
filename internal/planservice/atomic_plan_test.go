package planservice_test

import (
	"errors"
	"testing"

	"recommendation/internal/planbuilder"
	"recommendation/internal/plancache"
	"recommendation/internal/planservice"
)

func TestFailedPlanBuildDoesNotPublishPartialState(t *testing.T) {
	cache := plancache.New()
	svc := planservice.New(cache)
	p, err := svc.Generate("bad-plan", "panic")
	if !errors.Is(err, planbuilder.ErrBuild) || p != nil {
		t.Errorf("failed build returned plan=%v err=%v, want nil plan and build error", p, err)
	}
	if _, ok := cache.Get("bad-plan"); ok {
		t.Errorf("failed plan remained visible in cache")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("reading failed plan panicked: %v", r)
			}
		}()
		if _, ok := svc.StageName("bad-plan"); ok {
			t.Errorf("failed plan has a stage")
		}
	}()
	good, err := svc.Generate("good-plan", "rank")
	if err != nil || good == nil {
		t.Fatalf("unrelated valid plan failed after recovery: plan=%v err=%v", good, err)
	}
	if name, ok := svc.StageName("good-plan"); !ok || name != "rank" {
		t.Fatalf("valid plan stage=(%q,%v), want rank,true", name, ok)
	}
}
