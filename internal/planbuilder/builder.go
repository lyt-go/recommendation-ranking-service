package planbuilder

import "errors"

var ErrBuild = errors.New("ranking plan build failed")

type Stage struct{ Name string }
type Plan struct {
	Key      string
	Stage    *Stage
	Features map[string]float64
}

func Populate(plan *Plan, stage string) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrBuild
		}
	}()
	plan.Features["freshness"] = 1
	if stage == "panic" {
		panic("invalid stage")
	}
	plan.Stage = &Stage{Name: stage}
	return nil
}
