package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"stream-video/pkg/pb"
	"stream-video/pkg/repo/interfaces"

	"github.com/google/uuid"
)

const storageLocation = "storage"

type VideoServer struct {
	Repo interfaces.VideoRepo
	pb.UnimplementedVideoServiceServer
} 

func NewVideServer(repo interfaces.VideoRepo) pb.VideoServiceServer {
	return &VideoServer{
		Repo: repo,
	}
}

func (c *VideoServer) UploadVideo(stream pb.VideoService_UploadVideoServer) error {
	fileuid := uuid.New()
	fileName := fileuid.String()
	folderpath := storageLocation + "/" + fileName
	filepath := folderpath + "/" + fileName + ".mp4"
	if err := os.MkdirAll(folderpath, 0755); err != nil {
		return errors.New("failed to create directory")
	}
	newfile, err1 := os.Create(filepath)
	if err1 != nil {
		return errors.New("failed to create file")
	}
	defer newfile.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, err := newfile.Write(chunk.Data); err != nil {
			return err
		}
	}
	chanerr := make(chan error, 2)

	go func() {
		err := CreatePlaylistAndSegments(filepath, folderpath)
		chanerr <- err
	}()
	go func() {
		err := c.Repo.CreateVideo(fileName)
		chanerr <- err
	}()

	for i := 1; i <= 2; i++ {
		err := <-chanerr
		if err != nil {
			return err
		}
	}
	return stream.SendAndClose(&pb.UploadVideoResponse{
		Status:  http.StatusOK,
		Message: "Video succesfully uploaded",
		VideoId: fileName,
	})
}

func (c *VideoServer) StreamVideo(req *pb.StreamVideoRequest, stream pb.VideoService_StreamVideoServer) error {
	chunkSize := 4096
	buffer := make([]byte, chunkSize)
	playlistPath := fmt.Sprintf("storage/%s/%s", req.Videoid, req.Playlist)
	playlistFile, _ := os.Open(playlistPath)
	defer playlistFile.Close()
	for {
		n, err := playlistFile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&pb.StreamVideoResponse{
			VideoChunk: buffer[:n],
		}); err != nil {
			return err
		}

	}
	return nil
}

func (c *VideoServer) FindAllVideo(ctx context.Context, req *pb.FindAllRequest) (*pb.FindAllResponse, error) {
	res, err := c.Repo.FindAllVideo()
	if err != nil {
		return nil, err
	}
	return &pb.FindAllResponse{
		Status: http.StatusOK,
		Videos: res,
	}, nil
}

func CreatePlaylistAndSegments(filepath string, folderpath string) error {
	segmentDuration := 3
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", filepath,
		"-profile:v", "baseline", 
		"-level", "3.0",
		"-start_number", "0", 
		"-hls_time", strconv.Itoa(segmentDuration), 
		"-hls_list_size", "0", 
		"-f", "hls",
		fmt.Sprintf("%s/playlist.m3u8", folderpath),
	)
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v \nOutput: %s ", err, string(output))
	}
	return nil
}
