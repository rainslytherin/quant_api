package api

import (
	"fmt"

	"quant_api/config"
)

type Config struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Backend map[string]struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
}

func (c *Config) GetBackend(name string) string {
	if c.Backend == nil {
		return ""
	}
	if v, ok := c.Backend[name]; ok {
		return fmt.Sprintf("%s:%d", v.Host, v.Port)
	}
	return ""
}

func NewConfig(c *config.Config) *Config {
	return &Config{
		Host:    c.Http.Host,
		Port:    c.Http.Port,
		Backend: c.Backend,
	}
}
