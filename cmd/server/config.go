package main

import (
	"go.uber.org/fx"

	"go_sample_code/internal/database"
)

type ConfigPath string

type CommonConfig struct {
	Listen string `yaml:"listen"`
}

type Config struct {
	fx.Out
	CommonConfig   *CommonConfig            `yaml:"common"`
	DatabaseConfig *database.DatabaseConfig `yaml:"database"`
}
