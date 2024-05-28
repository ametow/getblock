package config

type Config struct {
	ApiKey         string
	GoroutineCount int     `yaml:"GoroutineCount"`
	BaseApiUrl     string  `yaml:"BaseApiUrl"`
	ListenAddress  string  `yaml:"ListenAddress"`
	EthPrice       float64 `yaml:"EthPrice"`
	BlockCount     int     `yaml:"BlockCount"`
}
