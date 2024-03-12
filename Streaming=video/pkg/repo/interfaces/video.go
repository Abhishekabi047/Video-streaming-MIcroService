package interfaces

import "stream-video/pkg/pb"

type VideoRepo interface {
	CreateVideo(string) error
	FindAllVideo() ([]*pb.VideoID, error)
}
