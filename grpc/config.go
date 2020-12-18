package grpc

type Config struct {
	Mode        string `mapstructure:"mode"`
	Port        string `mapstructure:"port"`
	RequestDump bool   `mapstructure:"request_dump"` // true or false
}
