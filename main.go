package main

import (
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
	fmt.Fprintf(os.Stdout, "configPath: %s\n", args.configPath)
}
