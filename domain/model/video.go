package model

import "time"

type Video struct {
	ID                     int64     `json:"id"`
	YoutubeVideoID         string    `json:"youtube_video_id"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	CreatedBy              int64     `json:"created_by"`
	UpdatedBy              int64     `json:"updated_by"`
	YoutubeTitle           string    `json:"youtube_title"`
	YoutubeDescription     string    `json:"youtube_description"`
	YoutubePlaylist        string    `json:"youtube_playlist"`
	YoutubeChannelID       string    `json:"youtube_channel_id"`
	YoutubeChannelUsername string    `json:"youtube_channel_username"`
	YoutubePrivacyStatus   string    `json:"youtube_privacy_status"`
}

// TableName sets the insert table name for this struct type
func (Video) TableName() string {
	return "video"
}
