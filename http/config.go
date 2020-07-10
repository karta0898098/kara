package http

type Config struct {
	Mode string `mapstructure:"mode"`
	Port string `mapstructure:"port"`
}

func (c *Config) New() Config {
	return *c
}
