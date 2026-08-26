// Package config 负责从环境变量加载服务配置。
package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Addr        string
	MaxPageSize int
	DefaultTopN int
}

func Load() *Config {
	cfg := &Config{
		Addr:        ":" + getenv("PORT", "8080"),
		MaxPageSize: getenvInt("MAX_PAGE_SIZE", 100),
		DefaultTopN: getenvInt("DEFAULT_TOP_N", 10),
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func (c *Config) String() string {
	return fmt.Sprintf("addr=%s max_page_size=%d default_top_n=%d", c.Addr, c.MaxPageSize, c.DefaultTopN)
}
