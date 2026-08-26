package ruleengine_test

import (
	"errors"
	"recommendation/internal/ruleconfig"
	"recommendation/internal/ruleengine"
	"testing"
)

func addRule(engine *ruleengine.Engine, name string, weight float64) (err error, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	return engine.Add(name, weight), false
}

func TestDefaultRuleConfigRejectsInvalidWeightAndRemainsWritable(t *testing.T) {
	engine := ruleengine.New(ruleconfig.Load(nil))
	err, panicked := addRule(engine, "freshness", -1)
	if panicked {
		t.Errorf("invalid default rule panicked")
	}
	if !errors.Is(err, ruleengine.ErrInvalidWeight) {
		t.Errorf("invalid weight error=%v, want validation error", err)
	}
	err, panicked = addRule(engine, "freshness", 1.5)
	if panicked {
		t.Errorf("valid rule after rejection panicked")
	} else if err != nil {
		t.Errorf("valid rule after rejection returned %v", err)
	}
	if value, ok := engine.Weight("freshness"); !ok || value != 1.5 {
		t.Errorf("stored weight=(%v,%v), want 1.5,true", value, ok)
	}
}
