package redis

type Config struct {
	Addrs    []string `yaml:"addrs"`
	Password string   `yaml:"password"`
	DB       int      `yaml:"db"`
	Master   string   `yaml:"master"`
}
