package http

type Config struct {
	Mode string `mapstructure:"mode"`
	Port string `mapstructure:"port"`
	Dump bool   `mapstructure:"dump"`
}

