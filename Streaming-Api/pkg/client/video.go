package client

import (
	"context"
	"io"
	"mime/multipart"
	"stream/pkg/client/interfaces"
	"stream/pkg/config"
	"stream/pkg/pb"

	errr "github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type VideoClient struct {
	Server pb.VideoServiceClient
}

func InitClient(c *config.Config) (pb.VideoServiceClient, error) {
	cc, err := grpc.Dial(c.VideoService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return pb.NewVideoServiceClient(cc), nil
}

func NewVideoClient(server pb.VideoServiceClient) interfaces.VideoClient {
	return &VideoClient{
		Server: server,
	}
}

func (c *VideoClient) UploadVideo(ctx context.Context, file *multipart.FileHeader) (*pb.UploadVideoResponse, error) {
	uploadFile, err := file.Open()
	if err != nil {
		return nil, err
	}

	defer uploadFile.Close()
	stream, err := c.Server.UploadVideo(ctx)
	if err != nil {
		return nil, errr.Wrap(err, "failed to start upload stream")
	}
	chunkSize := 4096
	buffer := make([]byte, chunkSize)
	for {
		n, err := uploadFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if err := stream.Send(&pb.UploadVideoRequest{
			Filename: file.Filename,
			Data:     buffer[:n],
		}); err != nil {
			return nil, err
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return response, nil

}

func (c *VideoClient) StreamVideo(ctx context.Context, filename, playlist string) (pb.VideoService_StreamVideoClient, error) {
	res, err := c.Server.StreamVideo(ctx, &pb.StreamVideoRequest{
		Videoid:  filename,
		Playlist: playlist,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *VideoClient) FindAllVideo(ctx context.Context) (*pb.FindAllResponse, error) {
	res, err := c.Server.FindAllVideo(ctx, &pb.FindAllRequest{})
	if err != nil {
		return nil, err
	}
	return res, nil
}
