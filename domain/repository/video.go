package repository

import "my-project/domain/model"

type IVideo interface {
	GetVideos() (*[]model.Video, error)
	InsertVideo(video *model.Video) error
	GetVideoByVideoID(videoID string) (*model.Video, error)
}
