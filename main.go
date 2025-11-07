package main

import (
	config "azarole/internal"
	"flag"
	"log/slog"
	"os"
)

type arguments struct {
	configPath string
}

func newArguments() *arguments {
	var configPath string
	flag.StringVar(&configPath, "config-path", "azarole.toml", "specify path to configuration file")
	flag.Parse()

	return &arguments{configPath: configPath}
}

func setupLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}

func main() {
	args := newArguments()
	config, err := config.LoadConfig(args.configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(111)
	}

	setupLogger()
	slog.Debug("Config loaded", "app.BaseUrl", config.App.BaseUrl)
}
