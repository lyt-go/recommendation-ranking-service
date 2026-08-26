package ruleconfig

type Config struct {
	Weights map[string]float64
	Enforce bool
}

func Load(values map[string]float64) *Config {
	if values == nil {
		return &Config{}
	}
	return &Config{Weights: values, Enforce: true}
}
