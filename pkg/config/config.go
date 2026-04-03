package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Since  string   `yaml:"since"`
	Until  string   `yaml:"until"`
	Repos  []string `yaml:"repos"`
	Output string   `yaml:"output"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Output == "" {
		cfg.Output = "report.html"
	}
	if cfg.Until == "" {
		cfg.Until = time.Now().Format("2006-01-02")
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Since == "" {
		return fmt.Errorf("since is required")
	}
	if _, err := time.Parse("2006-01-02", c.Since); err != nil {
		return fmt.Errorf("invalid since date %q: %w", c.Since, err)
	}
	if _, err := time.Parse("2006-01-02", c.Until); err != nil {
		return fmt.Errorf("invalid until date %q: %w", c.Until, err)
	}
	if len(c.Repos) == 0 {
		return fmt.Errorf("at least one repo is required")
	}
	for _, repo := range c.Repos {
		gitDir := filepath.Join(repo, ".git")
		if info, err := os.Stat(gitDir); err != nil || !info.IsDir() {
			return fmt.Errorf("not a git repository: %s", repo)
		}
	}
	return nil
}

func (c *Config) ParsedSince() (time.Time, error) {
	return time.Parse("2006-01-02", c.Since)
}

func (c *Config) ParsedUntil() (time.Time, error) {
	return time.Parse("2006-01-02", c.Until)
}
