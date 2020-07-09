package zlog

type Config struct {
	Env   string
	AppID string
	Debug bool
	Local bool
}

func (c *Config) New() Config {
	return *c
}
