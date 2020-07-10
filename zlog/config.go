package zlog

type Config struct {
	Env   string `mapstructure:"env"`
	AppID string `mapstructure:"app_id"`
	Debug bool   `mapstructure:"debug"`
	Local bool   `mapstructure:"local"`
}

func (c *Config) New() Config {
	return *c
}
