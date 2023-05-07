package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func Load() (Config, error) {
	var config Config
	b, err := os.ReadFile("config/application.yaml")
	if err != nil {
		return config, fmt.Errorf("error reading config file, %s", err)
	}

	if err := yaml.Unmarshal(b, &config); err != nil {
		return config, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return config, nil
}

type Config struct {
	Server struct {
		Addr        string   `yaml:"addr"`
		Debug       bool     `yaml:"debug"`
		SwaggerUi   bool     `yaml:"swagger-ui"`
		OpenapiSpec []string `yaml:"openapi-spec"`
	} `yaml:"server"`
}
