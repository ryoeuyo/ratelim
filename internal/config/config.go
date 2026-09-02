package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server    Server              `yaml:"server"`
	Log       Log                 `yaml:"log"`
	Storage   Storage             `yaml:"storage"`
	Limits    map[string]Limit    `yaml:"limits"`
	Upstreams map[string]Upstream `yaml:"upstreams"`
	Routes    []Route             `yaml:"routes"`
}

type Server struct {
	Addr            string        `yaml:"addr" env:"SERVER_ADDR" env-default:"localhost:8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" env-default:"15s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
}

type Storage struct {
	Driver string `yaml:"driver" env:"STORAGE_DRIVER" env-default:"redis"`
	Redis  Redis  `yaml:"redis"`
}

type Redis struct {
	Addr      string `yaml:"addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	DB        int    `yaml:"db" env:"REDIS_DB" env-default:"0"`
	Password  string `yaml:"password" env:"REDIS_PASSWORD"`
	KeyPrefix string `yaml:"key_prefix" env:"REDIS_KEY_PREFIX" env-default:"ratelim"`
}

type Limit struct {
	Rate   int           `yaml:"rate"`
	Window time.Duration `yaml:"window"`
	Key    []string      `yaml:"key"`
}

type Upstream struct {
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

type Route struct {
	Match    Match   `yaml:"match"`
	Upstream string  `yaml:"upstream"`
	Rewrite  Rewrite `yaml:"rewrite"`
	Limit    string  `yaml:"limit"`
}

type Match struct {
	PathPrefix string `yaml:"path_prefix"`
}

type Rewrite struct {
	StripPrefix string `yaml:"strip_prefix"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}
