package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ikitiki/logical_backup/pkg/config"
	"github.com/ikitiki/logical_backup/pkg/logger"
	"github.com/ikitiki/logical_backup/pkg/logicalbackup"
)

var (
	configFile  = flag.String("config", "config.yaml", "path to the config file")
	version     = flag.Bool("version", false, "print version information")
	development = flag.Bool("log-development", false, "enable development logging mode")

	Version  string
	Revision string

	GoVersion = runtime.Version()
)

func buildInfo() string {
	return fmt.Sprintf("logical backup version %s git revision %s go version %s", Version, Revision, GoVersion)
}

func main() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", buildInfo())
		_, _ = fmt.Fprintf(os.Stderr, "\nUsage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if *version {
		fmt.Println(buildInfo())
		os.Exit(1)
	}

	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		_, _ = fmt.Fprintf(os.Stderr, "Config file %s does not exist", *configFile)
		os.Exit(1)
	}

	cfg, err := config.New(*configFile, *development)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Could not load config file: %v", err)
		os.Exit(1)
	}
	// Initialize the logger after we've resolved the debug flag but before its first usage at cfg.Print
	if err := logger.InitGlobalLogger(cfg.Log); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Could not initialize global logger")
		os.Exit(1)
	}

	cfg.Print()

	lb, err := logicalbackup.New(cfg)
	if err != nil {
		logger.G.Fatalf("could not create backup instance: %v", err)
	}

	if err := lb.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
