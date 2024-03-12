package api

import (
	"fmt"
	"net"
	"stream-video/pkg/config"
	"stream-video/pkg/pb"

	"google.golang.org/grpc"
)

type Server struct {
	gs   *grpc.Server
	Lis  net.Listener
	Port string
}

func NewGrpcServe(c *config.Config, service pb.VideoServiceServer) (*Server, error) {
	grpcserver := grpc.NewServer()
	pb.RegisterVideoServiceServer(grpcserver, service)
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		return nil, err
	}
	return &Server{
		gs:   grpcserver,
		Lis:  lis,
		Port: c.Port,
	}, nil
}

func (s *Server) Start() error {
	fmt.Println("Video service on:", s.Port)
	return s.gs.Serve(s.Lis)
}
