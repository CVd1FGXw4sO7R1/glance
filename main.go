package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	app, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	log.Printf("starting glance v%s on %s:%d", version, cfg.Server.Host, cfg.Server.Port)

	if err := app.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
