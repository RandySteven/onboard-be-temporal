package configs

import (
	"flag"
	"fmt"
	"os"
)

func ParseFlags() (string, error) {
	var configPath string
	if flag.Lookup("config") == nil {
		flag.StringVar(&configPath, "config", "./files/yaml/app.local.yml", "path to config file")
		flag.Parse()
	} else {
		configPath = flag.Lookup("config").Value.String()
	}
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}
	return configPath, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("invalid file format")
	}
	return nil
}
