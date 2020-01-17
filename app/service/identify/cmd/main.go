package main

import (
	"github.com/liangjfblue/gmicro/app/service/identify/server"
)

const (
	webSrvName    = "gmicro.srv.identify"
	webSrvVersion = "v1.0.0"
)

func main() {
	srv := server.NewServer(webSrvName, webSrvVersion)
	srv.Init()

	srv.Run()
}
