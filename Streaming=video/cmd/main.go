package main

import (
	"log"
	"stream-video/pkg/config"
	"stream-video/pkg/di"
)

func main() {
	c,err:=config.LoadConfig()
	if err != nil {
		log.Fatal("failed to laod config",err.Error())
	}
	server,err1:=di.InitializeServe(c)
	if err1 != nil{
		log.Fatalf("failed to init server",err1.Error())
	}
	if err :=server.Start();err != nil{
		log.Fatal("Coundnt start the server")
	}
}