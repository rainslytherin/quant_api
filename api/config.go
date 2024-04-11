package api

import (
	"quant_api/config"
)

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewConfig(c *config.Config) *Config {
	return &Config{
		Host: c.Http.Host,
		Port: c.Http.Port,
	}
}
