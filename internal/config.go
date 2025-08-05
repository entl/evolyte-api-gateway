package internal

import (
	"io"
	"log"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Gateway  GatewayConfig
	JWT      JWTConfig
	Postgres Postgres
	Redis    Redis
}

type RouteConfig struct {
	Method     string `yaml:"method"`
	PathPrefix string `yaml:"path_prefix"`
}

type ServiceConfig struct {
	Name         string        `yaml:"name"`
	Backend      string        `yaml:"backend"`
	PublicRoutes []RouteConfig `yaml:"public_routes"`
}

type GatewayConfig struct {
	Services []ServiceConfig `yaml:"services"`
}

type JWTConfig struct {
	JWTSecret                     string `env:"JWT_SECRET_KEY,required"`
	JWTAlgorithm                  string `env:"JWT_ALGORITHM,required"`
	JWTTokenExpirationTime        int64  `env:"JWT_TOKEN_EXPIRATION_TIME,required"`
	JWTRefreshTokenExpirationTime int64  `env:"JWT_REFRESH_TOKEN_EXPIRATION_TIME,required"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     int    `env:"POSTGRES_PORT,required"`
	Username string `env:"POSTGRES_USERNAME,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Database string `env:"POSTGRES_DATABASE,required"`
}

type Redis struct {
	Host        string `env:"REDIS_HOST,required"`
	Port        int    `env:"REDIS_PORT,required"`
	Password    string `env:"REDIS_PASSWORD,required"`
	CacheDB     int    `env:"REDIS_CACHE_DB,required"`
	RateLimitDB int    `env:"REDIS_RATE_LIMIT_DB,required"`
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
