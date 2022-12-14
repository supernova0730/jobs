package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"net"
	"time"
)

type Config struct {
	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     string `envconfig:"DB_PORT" required:"true"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	DBName     string `envconfig:"DB_NAME" required:"true"`
	DBSSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`

	SchedulerRefreshRate time.Duration `envconfig:"SCHD_REFRESH_RATE" default:"30s"`

	HTTPHost string `envconfig:"HTTP_HOST" default:"localhost"`
	HTTPPort string `envconfig:"HTTP_PORT" required:"true"`
}

func Load(envPrefix string, filenames ...string) (Config, error) {
	config := Config{}

	err := godotenv.Load(filenames...)
	if err != nil {
		return Config{}, err
	}

	err = envconfig.Process(envPrefix, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) HttpAddress() string {
	return net.JoinHostPort(c.HTTPHost, c.HTTPPort)
}
