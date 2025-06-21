package config

import "time"

// app config

type App struct {
	Name     string        `mapstructure:"name"`
	Port     int           `mapstructure:"port"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

// store config

type Postgres struct {
	Driver string `mapstructure:"driver"`
	Dsn    string `mapstructure:"dsn"`
}

type Store struct {
	Postgres Postgres `mapstructure:"postgres"`
}
