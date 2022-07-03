package rest

import (
	"fmt"
	"github.com/cocoup/go-smart/core/prometheus"
	"github.com/cocoup/go-smart/rest/middleware"
	"github.com/cocoup/go-smart/tools/gocli/utils"
	"github.com/gin-gonic/gin"
	"log"
)

type (
	RunOption func(*Server)

	Server struct {
		Conf   RestConf
		Engine *gin.Engine
	}
)

func WithCors() RunOption {
	return func(server *Server) {
		server.Engine.Use(middleware.Cors())
	}
}

func MustNewServer(conf RestConf, opts ...RunOption) *Server {
	server, err := NewServer(conf, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return server
}

func NewServer(conf RestConf, opts ...RunOption) (*Server, error) {
	server := &Server{
		Conf:   conf,
		Engine: gin.New(),
	}

	server.Engine.ContextWithFallback = true

	if conf.CorsEnable {
		server.Engine.Use(middleware.Cors())
	}
	promMonitor := prometheus.NewMonitor(utils.ToSnakeCase(conf.Name))
	server.Engine.Use(middleware.Prometheus(promMonitor, server.Engine))
	server.Engine.Use(middleware.LogHandler(conf.Verbose))
	server.Engine.Use(middleware.Recovery(true))
	for _, opt := range opts {
		opt(server)
	}
	server.Engine.Use(middleware.Unzip())

	//server.Engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return server, nil
}

func (s *Server) Start() {
	if nil != s {
		s.Engine.Run(fmt.Sprintf("%s:%d", s.Conf.Host, s.Conf.Port))
	}
}
