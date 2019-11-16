package config

import (
	"fmt"
	"github.com/olebedev/config"
	"os"
)

var (
	Config *config.Config
)

func LoadCofig(path string) error {
	mode := os.Getenv("MODE")
	if mode == "" {
		mode = "debug"
	}

	cfg, err := config.ParseYamlFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse config from %s failed: %s", path, err)
		return err
	}

	if Config, err = cfg.Get(mode); err != nil {
		return err
	}
	return nil
}
