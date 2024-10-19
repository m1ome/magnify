package pkg

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Prometheus  ConfigPrometheus   `yaml:"prometheus"`
		Expressions []ConfigExperssion `yaml:"expressions"`
		Cache       ConfigCache
		Http        ConfigHttp `yaml:"http"`
	}

	ConfigExperssion struct {
		Name       string `yaml:"name"`
		Query      string `yaml:"query"`
		Experssion string `yaml:"expr"`
	}

	ConfigPrometheus struct {
		Addr string `yaml:"addr"`
	}

	ConfigCache struct {
		Expiration int64 `yaml:"expire"`
	}

	ConfigHttp struct {
		Addr string `yaml:"addr"`
	}
)

func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}

	cfg := &Config{}
	cfg.Cache.Expiration = 60
	cfg.Http.Addr = ":8080"

	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	return cfg, nil
}
