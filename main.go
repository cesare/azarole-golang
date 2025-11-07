package main

import (
	config "azarole/internal"
	"flag"
	"fmt"
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

func main() {
	args := newArguments()
	config, err := config.LoadConfig(args.configPath)
	if err != nil {
		os.Exit(111)
	}

	fmt.Fprintf(os.Stdout, "config.app.BaseUrl: %s\n", config.App.BaseUrl)
}
