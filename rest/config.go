package rest

import (
	"github.com/cocoup/go-smart/core/prometheus"
	"github.com/cocoup/go-smart/core/stores/redis"
	"github.com/cocoup/go-smart/core/stores/sqlx"
	"github.com/cocoup/go-smart/core/trace"
)

type ServiceConf struct {
	Name       string            `yaml:"name"`
	Mode       string            `yaml:"mode"`
	MetricsUrl string            `yaml:"metrics-url"`
	JWT        JWTConfig         `yaml:"jwt"`
	SqlConf    sqlx.Config       `yaml:"sql"`
	RedisConf  redis.Config      `yaml:"redis"`
	Prometheus prometheus.Config `yaml:"prometheus"`
	Trace      trace.Config      `yaml:"trace"`

	//Log        logx.LogConf
}

type RestConf struct {
	ServiceConf `yaml:",inline"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	RouteRoot   string `yaml:"route-root"`
	CorsEnable  bool   `yaml:"cors-enable"`
	Verbose     bool   `yaml:"verbose"` //日志内容(默认简单,设置true为详细日志)
}

type JWTConfig struct {
	Secret string `yaml:"secret"` // jwt密钥
	Exp    int64  `yaml:"exp"`    // 过期时间(秒)
	Issuer string `yaml:"issuer"` // 签发者
}
