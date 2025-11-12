package main

import (
	"azarole/internal/core"
	"azarole/internal/server"
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
	application, err := core.LoadApplication(args.configPath)
	if err != nil {
		slog.Error("Failed to load application", "error", err)
		os.Exit(111)
	}

	setupLogger()

	engine := server.Engine(application)
	engine.Run(application.Config.Server.BindAddress())
}
