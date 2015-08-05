package main

import (
	"github.com/grayzone/example/websrv"
	"log"
)

func main() {
	var srv = websrv.SrvOp{":12345", "/tmpl/"}
	log.Println(srv)
	srv.InitRoute()
	srv.Start()
}
