package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen          string        `yaml:"listen"`
	DBPath          string        `yaml:"db_path"`
	CollectInterval time.Duration `yaml:"collect_interval"`
	JWTSecret       string        `yaml:"jwt_secret"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"` // plain text
	Alert           AlertConfig   `yaml:"alert"`
}

type AlertConfig struct {
	CPU     float64 `yaml:"cpu"`
	Memory  float64 `yaml:"memory"`
	Disk    float64 `yaml:"disk"`
	Webhook string  `yaml:"webhook"`
}

func Default() *Config {
	return &Config{
		Listen:          "0.0.0.0:1080",
		DBPath:          "gopanel.db",
		CollectInterval: 5 * time.Second,
		JWTSecret:       "gopanel-change-me",
		Username:        "admin",
		Password:        "admin",
		Alert: AlertConfig{CPU: 90, Memory: 90, Disk: 90},
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := Default()
	return cfg, yaml.Unmarshal(data, cfg)
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
