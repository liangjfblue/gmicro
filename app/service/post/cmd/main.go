package main

import (
	"github.com/liangjfblue/gmicro/app/service/post/server"
)

const (
	webUserSrvName    = "gmicro.srv.post"
	webUserSrvVersion = "v1.0.0"
)

func main() {
	srv := server.NewServer(webUserSrvName, webUserSrvVersion)
	srv.Init()

	srv.Run()
}
