package di

import (
	"stream-video/pkg/api"
	"stream-video/pkg/api/service"
	"stream-video/pkg/config"
	"stream-video/pkg/db"
	"stream-video/pkg/repo"

	"github.com/google/wire"
)

func InitializeServe1(c *config.Config) (*api.Server, error) {
	wire.Build(db.InitDB, repo.NewVideoRepo, service.NewVideServer, api.NewGrpcServe)
	return &api.Server{}, nil
}
