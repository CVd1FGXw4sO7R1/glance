package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/glanceapp/glance/internal/config"
	"github.com/glanceapp/glance/internal/server"
)

const version = "0.1.0"

func main() {
	var (
		configPath  = flag.String("config", "glance.yml", "Path to the configuration file")
		showVersion = flag.Bool("version", false, "Print version and exit")
		host        = flag.String("host", "", "Host to listen on (overrides config)")
		port        = flag.Int("port", 0, "Port to listen on (overrides config)")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("glance v%s\n", version)
		os.Exit(0)
	}

	// Default config search order:
	// 1. Explicit -config flag value
	// 2. ./glance.yml (current directory)
	// 3. ~/.config/glance/glance.yml (XDG-style config dir)
	// 4. ~/glance.yml (home directory fallback)
	if _, err := os.Stat(*configPath); os.IsNotExist(err) && *configPath == "glance.yml" {
		if home, err := os.UserHomeDir(); err == nil {
			xdgPath := filepath.Join(home, ".config", "glance", "glance.yml")
			if _, err := os.Stat(xdgPath); err == nil {
				*configPath = xdgPath
			} else {
				*configPath = filepath.Join(home, "glance.yml")
			}
		}
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// CLI flags override config file values
	if *host != "" {
		cfg.Server.Host = *host
	}
	if *port != 0 {
		cfg.Server.Port = *port
	}

	// Default to localhost if no host is set, to avoid accidentally exposing
	// the server on all interfaces when running locally.
	if cfg.Server.Host == "" {
		cfg.Server.Host = "127.0.0.1"
	}

	// Default to port 3000; 8080 is frequently occupied by other dev servers.
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 3000
	}

	app, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	log.Printf("starting glance v%s on %s:%d", version, cfg.Server.Host, cfg.Server.Port)

	if err := app.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
