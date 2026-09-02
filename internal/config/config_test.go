package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	cfg, err := Load("../../config.yaml")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Addr != "localhost:8080" {
		t.Fatalf("Server.Addr = %q", cfg.Server.Addr)
	}
	if cfg.Server.ShutdownTimeout != 10*time.Second {
		t.Fatalf("Server.ShutdownTimeout = %s", cfg.Server.ShutdownTimeout)
	}
	if cfg.Storage.Redis.Addr != "localhost:6379" {
		t.Fatalf("Redis.Addr = %q", cfg.Storage.Redis.Addr)
	}
	if _, ok := cfg.Limits["default"]; !ok {
		t.Fatal("limits.default is missing")
	}
	if _, ok := cfg.Upstreams["users"]; !ok {
		t.Fatal("upstreams.users is missing")
	}
	if len(cfg.Routes) != 2 {
		t.Fatalf("len(Routes) = %d", len(cfg.Routes))
	}
	if cfg.Routes[0].Rewrite.StripPrefix != "/api" {
		t.Fatalf("Routes[0].Rewrite.StripPrefix = %q", cfg.Routes[0].Rewrite.StripPrefix)
	}
}
