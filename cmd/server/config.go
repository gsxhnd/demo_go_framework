package main

import "go.uber.org/fx"

type ConfigPath string

type CommonConfig struct {
	Listen string `yaml:"listen"`
}

type Config struct {
	fx.Out
	CommonConfig *CommonConfig `yaml:"common"`
}
