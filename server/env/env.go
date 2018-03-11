package env

import (
	"github.com/caarlos0/env"
)

const (
	DB_DRIVER = "postgres"
)

type Config struct {
	DB DB
}

type DB struct {
	Host string `env:"DBHOST" envDefault:"localhost"`
	Username string `env:"DBUSER" envDefault:"loganw"`
	Password string `env:"DBPASS" envDefault:"\"\""`
	Name string `env:"DBNAME" envDefault:"mta"`
}

func NewConfig() *Config {
	conf := &Config{}

	err := env.Parse(&conf.DB)
	if err != nil {
		panic(err)
	}

	return conf
}