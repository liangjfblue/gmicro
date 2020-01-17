package main

import (
	"github.com/liangjfblue/gmicro/app/interface/user/server"
)

const (
	webUserSrvName    = "gmicro.web.user"
	webUserSrvVersion = "v1.0.0"
)

func main() {
	srv := server.NewServer(webUserSrvName, webUserSrvVersion)
	srv.Init()

	srv.Run()
}
