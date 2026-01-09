package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}
type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	return MustLoadByPath(path)
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(fmt.Sprint("config file do not exists", path))
	}
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic(fmt.Sprint("faild to read config", err.Error()))
	}
	return &config
}

func fetchConfigPath() string {
	var pathConfig string
	// --path_config="../../config/local.yaml"
	flag.StringVar(&pathConfig, "config", "", "path to config file")
	flag.Parse()

	if pathConfig == "" {
		pathConfig = os.Getenv("CONFIG_PATH")
	}

	return pathConfig
}
