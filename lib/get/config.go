package get

type Config struct {
	Interval int `mapstructure:"interval"`
	Limit    int `mapstructure:"limit"`
}
