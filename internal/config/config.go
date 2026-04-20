package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Pages    []PageConfig   `yaml:"pages"`
	Branding BrandingConfig `yaml:"branding"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host       string        `yaml:"host"`
	Port       uint16        `yaml:"port"`
	AssetsPath string        `yaml:"assets-path"`
	BaseURL    string        `yaml:"base-url"`
	Timeout    time.Duration `yaml:"timeout"`
}

// BrandingConfig holds UI branding/customization settings.
type BrandingConfig struct {
	HideFooter      bool   `yaml:"hide-footer"`
	CustomFooterText string `yaml:"custom-footer-text"`
	LogoURL         string `yaml:"logo-url"`
	FaviconURL      string `yaml:"favicon-url"`
	AppName         string `yaml:"app-name"`
}

// PageConfig represents a single dashboard page.
type PageConfig struct {
	Name    string         `yaml:"name"`
	Slug    string         `yaml:"slug"`
	Columns []ColumnConfig `yaml:"columns"`
	HideDesktopNav bool   `yaml:"hide-desktop-nav"`
	ExpandMobilePageOnLoad bool `yaml:"expand-mobile-page-on-load"`
}

// ColumnConfig represents a column within a page.
type ColumnConfig struct {
	Size    string        `yaml:"size"`
	Widgets []WidgetConfig `yaml:"widgets"`
}

// WidgetConfig holds the raw configuration for any widget type.
// The Type field determines how the rest of the fields are interpreted.
type WidgetConfig struct {
	Type  string                 `yaml:"type"`
	Extra map[string]interface{} `yaml:",inline"`
}

// defaults applies default values to the configuration where fields are zero.
func (c *Config) defaults() {
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.Timeout == 0 {
		c.Server.Timeout = 5 * time.Second
	}
	if c.Branding.AppName == "" {
		c.Branding.AppName = "Glance"
	}
}

// Validate checks the configuration for required fields and logical errors.
func (c *Config) Validate() error {
	if len(c.Pages) == 0 {
		return fmt.Errorf("at least one page must be defined")
	}
	for i, page := range c.Pages {
		if page.Name == "" {
			return fmt.Errorf("page at index %d is missing a name", i)
		}
	}
	return nil
}

// Load reads and parses a YAML configuration file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	cfg.defaults()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}
