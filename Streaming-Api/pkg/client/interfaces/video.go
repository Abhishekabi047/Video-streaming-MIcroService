package interfaces

import (
	"context"
	"mime/multipart"
	"stream/pkg/pb"
)

type VideoClient interface {
	UploadVideo(context.Context, *multipart.FileHeader) (*pb.UploadVideoResponse, error)
	StreamVideo(context.Context, string, string) (pb.VideoService_StreamVideoClient, error)
	FindAllVideo(context.Context) (*pb.FindAllResponse, error)
}
