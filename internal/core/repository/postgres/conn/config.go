package core_postgres_pool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string        `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     string        `envconfig:"POSTGRES_PORT" default:"5434"`
	User     string        `envconfig:"POSTGRES_USER" required:"true"`
	Password string        `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Database string        `envconfig:"POSTGRES_DB" required:"true"`
	Timeout  time.Duration `envconfig:"POSTGRES_TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Postgres connection pool config: %w", err)
		panic(err)
	}
	return config
}
