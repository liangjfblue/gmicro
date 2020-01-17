package main

import (
	"github.com/liangjfblue/gmicro/app/interface/comment/server"
)

const (
	webSrvName    = "gmicro.web.comment"
	webSrvVersion = "v1.0.0"
)

func main() {
	srv := server.NewServer(webSrvName, webSrvVersion)
	srv.Init()

	srv.Run()
}
