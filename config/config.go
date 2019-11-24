package config

import (
	"fmt"
	"github.com/olebedev/config"
	"os"
)

const CONFIG_PATH = "config/config.yml"

var (
	Config *config.Config
)

func LoadConfig(path ...string) error {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "debug"
	}
	filePath := CONFIG_PATH
	if len(path) > 0 {
		filePath = path[0]
	}
	cfg, err := config.ParseYamlFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse config from %s failed: %s", path, err)
		return err
	}

	if Config, err = cfg.Get(mode); err != nil {
		return err
	}
	return nil
}
