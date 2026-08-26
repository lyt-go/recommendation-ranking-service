package ruleconfig

type Config struct {
	Weights map[string]float64
	Enforce bool
}

func Load(values map[string]float64) *Config {
	if values == nil {
		// 未提供特征规则配置时使用默认配置：同样启用规则校验，
		// 并初始化可写的空权重表，避免写入时触发 nil map panic。
		return &Config{Weights: make(map[string]float64), Enforce: true}
	}
	return &Config{Weights: values, Enforce: true}
}
