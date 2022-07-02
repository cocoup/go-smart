package main

import (
	"flag"

	{{.imports}}
)

var configFile = flag.String("c", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

    svcCtx := service.NewContext(&c)

    server := rest.MustNewServer(
        c.RestConf,
        rest.WithRecovery(true),
    )

    route.RegisterHandlers(server, svcCtx)
    server.Start()
}
