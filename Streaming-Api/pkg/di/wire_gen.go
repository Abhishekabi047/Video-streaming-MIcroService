// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"stream/pkg/api"
	"stream/pkg/api/handler"
	"stream/pkg/client"
	"stream/pkg/config"
)

// Injectors from wire.go:

func InitializeAPI(c *config.Config) (*api.Server, error) {
	videoServiceClient, err := client.InitClient(c)
	if err != nil {
		return nil, err
	}
	videoClient := client.NewVideoClient(videoServiceClient)
	videoHandler := handler.NewVideoHandler(videoClient)
	server, err := api.NewServerHTTP(c, videoHandler)
	if err != nil {
		return nil, err
	}
	return server, nil
}
