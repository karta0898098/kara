package zlog

type Config struct {
	Env   string `mapstructure:"env"`
	AppID string `mapstructure:"app_id"`
	Debug bool   `mapstructure:"debug"`
	Level int8   `mapstructure:"level"`
	Local bool   `mapstructure:"local"`
}
