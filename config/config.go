package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

type Env string

const (
	Default Env = ""
	Dev     Env = "dev"
	Prod    Env = "prod"
)

type (
	LoadConfig struct {
		Env      Env
		DirPath  string
		FileName string
	}
	Config struct {
		EnvName string `yaml:"env-name"`
		Server  Server `yaml:"server"`
		DB      DB     `yaml:"db"`
	}
	Server struct {
		Addr        string        `yaml:"addr"`
		Timeout     time.Duration `yaml:"timeout"`
		Debug       bool          `yaml:"debug"`
		SwaggerUi   bool          `yaml:"swagger-ui"`
		OpenapiSpec []string      `yaml:"openapi-spec"`
	}
	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}
)

func NewLoadConfig(env string) LoadConfig {
	return LoadConfig{
		Env:      Env(env),
		DirPath:  "config",
		FileName: "application",
	}
}

func MustLoad(loadCfg LoadConfig) Config {
	cfg, err := Load(loadCfg)
	if err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}

	return cfg
}

// Load and return app configuration.
// Configuration is constructed from multiple sources in following order:
//  1. application.yaml - default config for all environments
//  2. application-${env-name}.yaml - environment specific config file
//  3. environment variables - only consider vars with FGO_ prefix. _ is path delimiter, while __ is used instead of
//     dashes. E.g. FGO_SERVER_SWAGGER__UI=false will affect a config with the key server.swagger-ui.
//
// All sources are merged at the end, where later sources have higher priority and will override previous ones.
func Load(loadCfg LoadConfig) (Config, error) {
	cfg := Config{}

	// default file
	if err := load(cfgFilePath(loadCfg.FileName, loadCfg.DirPath, Default)); err != nil {
		return cfg, err
	}

	// env specific file
	if len(loadCfg.Env) != 0 {
		if err := load(cfgFilePath(loadCfg.FileName, loadCfg.DirPath, loadCfg.Env)); err != nil {
			return cfg, err
		}
	}

	// env vars
	if err := k.Load(env.Provider("FGO_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "FGO_")
		s = strings.ToLower(s)
		s = strings.Replace(s, "__", "-", -1)
		return strings.Replace(s, "_", ".", -1)
	}), nil); err != nil {
		return cfg, fmt.Errorf("error loading env vars: %v", err)
	}

	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return cfg, fmt.Errorf("error unmarshaling to a struct: %v", err)
	}

	return cfg, nil
}

func load(path string) error {
	if err := k.Load(file.Provider(path), kyaml.Parser()); err != nil {
		return fmt.Errorf("error loading '%s' config file: %v", path, err)
	}
	return nil
}

func cfgFilePath(fileName, dirPath string, env Env) string {
	name := fileName
	if env != Default {
		name += "-" + string(env)
	}
	name += "." + "yaml"
	return filepath.Join(dirPath, name)
}
