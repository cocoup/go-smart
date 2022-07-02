package rest

import (
	"github.com/cocoup/go-smart/core/jwt"
	"github.com/cocoup/go-smart/core/prometheus"
	"github.com/cocoup/go-smart/core/trace"
)

type ServiceConf struct {
	Name       string `yaml:"name"`
	Mode       string `yaml:"mode"`
	MetricsUrl string `yaml:"metrics-url"`
	JWT        jwt.Config
	Prometheus prometheus.Config
	Trace      trace.Config
	//Log        logx.LogConf
}

type RestConf struct {
	ServiceConf `yaml:",inline"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	RouteRoot   string `yaml:"route-root"`
	CorsEnable  bool   `yaml:"cors-enable"`
}
