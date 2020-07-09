package http

type Config struct {
	Mode string
	Port string
}

func (c *Config) New() Config {
	return *c
}
