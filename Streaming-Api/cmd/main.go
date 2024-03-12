package main

import (
	"log"
	"stream/pkg/config"
	"stream/pkg/di"

)

func main() {
	c,configerr:=config.LoadConfig()
	if configerr != nil{
		log.Fatal("cannot load config",configerr)
	}
	server,dier:=di.InitializeAPI(c)
	if dier != nil{
		log.Fatal("cannot intitalize server",dier)
	}
	server.Start()
}