package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type Config struct {
	Gateway GatewayConfig
	JWT     JWTConfig
}

type RouteConfig struct {
	FromPath     string `yaml:"from_path"`
	ToPath       string `yaml:"to_path"`
	Method       string `yaml:"method"`
	AuthRequired bool   `yaml:"auth_required"`
}

type ServiceConfig struct {
	Backend      string        `yaml:"backend"`
	PublicPrefix string        `yaml:"public_prefix"`
	Routes       []RouteConfig `yaml:"routes"`
}

type GatewayConfig struct {
	Services map[string]ServiceConfig `yaml:"services"`
}

type JWTConfig struct {
	JWTSecret                     string `env:"JWT_SECRET_KEY,required"`
	JWTAlgorithm                  string `env:"JWT_ALGORITHM,required"`
	JWTTokenExpirationTime        int64  `env:"JWT_TOKEN_EXPIRATION_TIME,required"`
	JWTRefreshTokenExpirationTime int64  `env:"JWT_REFRESH_TOKEN_EXPIRATION_TIME,required"`
}

func LoadConfig(envFile string, proxyConfigFile string) (*Config, error) {
	var cfg Config
	_ = godotenv.Load(envFile)

	// Load proxy configuration
	proxyConfig, err := loadProxyConfig(proxyConfigFile)
	if err != nil {
		return nil, err
	}
	cfg.Gateway = *proxyConfig

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse env: %v", err)
	}

	return &cfg, nil
}

func loadProxyConfig(path string) (*GatewayConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body, err := io.ReadAll(f)

	var cfg GatewayConfig
	if err := yaml.Unmarshal(body, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
