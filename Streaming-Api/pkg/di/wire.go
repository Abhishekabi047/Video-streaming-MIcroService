package di

import (
	"stream/pkg/api"
	"stream/pkg/api/handler"
	"stream/pkg/client"
	"stream/pkg/config"

	"github.com/google/wire"
)

func InitializeAPI1(c *config.Config) (*api.Server,error) {
	wire.Build(client.InitClient,client.NewVideoClient,handler.NewVideoHandler,api.NewServerHTTP)
	return &api.Server{},nil
}