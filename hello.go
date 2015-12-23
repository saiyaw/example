package main

import (
	"log"

	"github.com/saiyawang/example/websrv"
)

func main() {
	var srv = websrv.SrvOp{":12345", "/tmpl/"}
	log.Println(srv)
	srv.InitRoute()
	srv.Start()
}
