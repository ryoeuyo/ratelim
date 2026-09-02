package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
)

var allowedKeyParts = map[string]struct{}{
	"ip":   {},
	"path": {},
}

func (c *Config) Validate() error {
	var errs []error

	if strings.TrimSpace(c.Server.Addr) == "" {
		errs = append(errs, errors.New("server.addr is required"))
	}
	if c.Server.ReadTimeout <= 0 {
		errs = append(errs, errors.New("server.read_timeout must be positive"))
	}
	if c.Server.WriteTimeout <= 0 {
		errs = append(errs, errors.New("server.write_timeout must be positive"))
	}
	if c.Server.ShutdownTimeout <= 0 {
		errs = append(errs, errors.New("server.shutdown_timeout must be positive"))
	}

	if _, err := parseLogLevel(c.Log.Level); err != nil {
		errs = append(errs, err)
	}

	if c.Storage.Driver != "redis" {
		errs = append(errs, fmt.Errorf("storage.driver %q is not supported", c.Storage.Driver))
	}
	if strings.TrimSpace(c.Storage.Redis.Addr) == "" {
		errs = append(errs, errors.New("storage.redis.addr is required"))
	}

	if len(c.Limits) == 0 {
		errs = append(errs, errors.New("limits must contain at least one policy"))
	}
	for name, limit := range c.Limits {
		if err := validateLimit(name, limit); err != nil {
			errs = append(errs, err)
			continue
		}
		if len(limit.Key) == 0 {
			limit.Key = []string{"ip", "path"}
			c.Limits[name] = limit
		}
	}

	if len(c.Upstreams) == 0 {
		errs = append(errs, errors.New("upstreams must contain at least one backend"))
	}
	for name, upstream := range c.Upstreams {
		if err := validateUpstream(name, upstream); err != nil {
			errs = append(errs, err)
		}
	}

	if len(c.Routes) == 0 {
		errs = append(errs, errors.New("routes must contain at least one route"))
	}
	for i, route := range c.Routes {
		if err := validateRoute(i, route, c.Limits, c.Upstreams); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func parseLogLevel(level string) (slog.Level, error) {
	var parsed slog.Level
	if err := parsed.UnmarshalText([]byte(level)); err != nil {
		return 0, fmt.Errorf("log.level %q is invalid", level)
	}

	return parsed, nil
}

func (c *Config) SlogLevel() slog.Level {
	level, _ := parseLogLevel(c.Log.Level)
	return level
}

func validateLimit(name string, limit Limit) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("limits: empty policy name")
	}
	if limit.Rate <= 0 {
		return fmt.Errorf("limits.%s.rate must be positive", name)
	}
	if limit.Window <= 0 {
		return fmt.Errorf("limits.%s.window must be positive", name)
	}
	for _, part := range limit.Key {
		if _, ok := allowedKeyParts[part]; !ok {
			return fmt.Errorf("limits.%s.key: unsupported part %q", name, part)
		}
	}

	return nil
}

func validateUpstream(name string, upstream Upstream) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("upstreams: empty backend name")
	}

	parsed, err := url.Parse(upstream.URL)
	if err != nil {
		return fmt.Errorf("upstreams.%s.url: %w", name, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("upstreams.%s.url must include scheme and host", name)
	}
	if upstream.Timeout <= 0 {
		return fmt.Errorf("upstreams.%s.timeout must be positive", name)
	}

	return nil
}

func validateRoute(i int, route Route, limits map[string]Limit, upstreams map[string]Upstream) error {
	prefix := route.Match.PathPrefix
	if !strings.HasPrefix(prefix, "/") {
		return fmt.Errorf("routes[%d].match.path_prefix must start with /", i)
	}
	if _, ok := upstreams[route.Upstream]; !ok {
		return fmt.Errorf("routes[%d].upstream %q is not defined", i, route.Upstream)
	}
	if _, ok := limits[route.Limit]; !ok {
		return fmt.Errorf("routes[%d].limit %q is not defined", i, route.Limit)
	}
	if strip := route.Rewrite.StripPrefix; strip != "" && !strings.HasPrefix(prefix, strip) {
		return fmt.Errorf("routes[%d].rewrite.strip_prefix %q is not a prefix of %q", i, strip, prefix)
	}

	return nil
}
