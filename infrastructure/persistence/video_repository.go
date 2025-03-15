package persistence

import (
	"gorm.io/gorm"
	"my-project/domain/model"
	"my-project/domain/repository"
	"my-project/infrastructure/logger"
)

type VideoRepository struct {
	DB *gorm.DB
}

func (v VideoRepository) GetVideoByVideoID(videoID string) (*model.Video, error) {
	if videoID == "" {
		return nil, nil
	}

	var video model.Video
	if err := v.DB.Where("youtube_video_id = ?", videoID).First(&video).Error; err != nil {
		logger.GetLogger().WithField("error", err).Error("Failed to get video by video id")
		return nil, err
	}

	return &video, nil
}

func (v VideoRepository) GetVideos() (*[]model.Video, error) {
	var videos []model.Video
	if err := v.DB.Find(&videos).Error; err != nil {
		logger.GetLogger().WithField("error", err).Error("Failed to get videos")
		return nil, err
	}

	return &videos, nil
}

func (v VideoRepository) InsertVideo(video *model.Video) error {
	if err := v.DB.Save(video).Error; err != nil {
		logger.GetLogger().WithField("error", err).Error("Failed to insert video")
		return err
	}

	return nil
}

func NewVideoRepository(db *gorm.DB) repository.IVideo {
	return &VideoRepository{DB: db}
}
