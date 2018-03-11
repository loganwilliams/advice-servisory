package env

import (
	"github.com/caarlos0/env"
	"log"
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
		log.Panic("Error parsing DB config environment variables", err)
	}

	return conf
}

type TestConfig struct {
	DB TestDB
}

type TestDB struct {
	Host string `env:"TEST_DBHOST" envDefault:"localhost"`
	Username string `env:"TEST_DBUSER" envDefault:"loganw"`
	Password string `env:"TEST_DBPASS" envDefault:"\"\""`
	Name string `env:"TEST_DBNAME" envDefault:"mta_test"`
}

func NewTestConfig() *TestConfig {
	conf := &TestConfig{}

	// parse server conf
	err := env.Parse(&conf.DB)
	if err != nil {
		log.Panic("Error parsing test DB config environment variables", err)
	}

	return conf
}