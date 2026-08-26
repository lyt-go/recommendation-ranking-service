package ruleengine

import (
	"errors"
	"recommendation/internal/ruleconfig"
)

var ErrInvalidWeight = errors.New("feature weight must be positive")

type Validator interface{ Validate(string, float64) error }
type weightValidator struct{}

func (v *weightValidator) Validate(_ string, weight float64) error {
	if weight <= 0 {
		return ErrInvalidWeight
	}
	return nil
}
func NewValidator(cfg *ruleconfig.Config) Validator {
	if !cfg.Enforce {
		// 未启用校验时显式返回 nil 接口；不要返回 typed-nil 指针，
		// 否则调用方 e.validator != nil 判定为 true，引发无效规则被静默放行。
		return nil
	}
	return &weightValidator{}
}

type Engine struct {
	cfg       *ruleconfig.Config
	validator Validator
}

func New(cfg *ruleconfig.Config) *Engine { return &Engine{cfg: cfg, validator: NewValidator(cfg)} }
func (e *Engine) Add(name string, weight float64) error {
	if e.validator != nil {
		if err := e.validator.Validate(name, weight); err != nil {
			return err
		}
	}
	e.cfg.Weights[name] = weight
	return nil
}
func (e *Engine) Weight(name string) (float64, bool) { v, ok := e.cfg.Weights[name]; return v, ok }
