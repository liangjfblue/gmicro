package main

import (
	"github.com/liangjfblue/gmicro/app/interface/post/server"
)

const (
	webPostSrvName    = "gmicro.web.post"
	webPostSrvVersion = "v1.0.0"
)

func main() {
	srv := server.NewServer(webPostSrvName, webPostSrvVersion)
	srv.Init()

	srv.Run()
}
