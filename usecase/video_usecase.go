package usecase

import (
	"flag"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"my-project/domain/model"
	"my-project/domain/repository"
	"my-project/infrastructure/logger"
	"time"
)

type IVideoUsecase interface {
	GetVideos() ([]model.Video, error)
}

type VideoUsecase struct {
	YoutubeService  *youtube.Service
	VideoRepository repository.IVideo
	playlistId      string
}

func NewVideoUsecase(service *youtube.Service, videoRepository repository.IVideo) IVideoUsecase {
	return &VideoUsecase{YoutubeService: service, VideoRepository: videoRepository}
}

var (
	method = flag.String("method", "list", "The API method to execute. (List is the only method that this sample currently supports.")

	channelId              = flag.String("channelId", "", "Retrieve playlists for this channel. Value is a YouTube channel ID.")
	hl                     = flag.String("hl", "", "Retrieve localized resource metadata for the specified application language.")
	maxResults             = flag.Int64("maxResults", 5, "The maximum number of playlist resources to include in the API response.")
	mine                   = flag.Bool("mine", false, "List playlists for authenticated user's channel. Default: false.")
	onBehalfOfContentOwner = flag.String("onBehalfOfContentOwner", "", "Indicates that the request's auth credentials identify a user authorized to act on behalf of the specified content owner.")
	pageToken              = flag.String("pageToken", "", "Token that identifies a specific page in the result set that should be returned.")
	part                   = flag.String("part", "snippet", "Comma-separated list of playlist resource parts that API response will include.")
	playlistId             = flag.String("playlistId", "", "Retrieve information about this playlist.")
)

// Retrieve playlistItems in the specified playlist
func (VideoUsecase *VideoUsecase) playlistItemsList(part []string, playlistId string, pageToken string) *youtube.PlaylistItemListResponse {
	call := VideoUsecase.YoutubeService.PlaylistItems.List(part)
	call = call.PlaylistId(playlistId)
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	response, err := call.Do()
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while getting playlist items")
	}
	return response
}

// Retrieve resource for the authenticated user's channel
func (VideoUsecase *VideoUsecase) channelsListMine(part []string) *youtube.ChannelListResponse {
	call := VideoUsecase.YoutubeService.Channels.List(part)
	call = call.Mine(true)
	response, err := call.Do()
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while getting channel list")
	}
	return response
}

func (VideoUsecase *VideoUsecase) playlistsList(part []string, channelId string, hl string, maxResults int64, mine bool, onBehalfOfContentOwner string, pageToken string) *youtube.PlaylistListResponse {
	call := VideoUsecase.YoutubeService.Playlists.List(part)
	if channelId != "" {
		call = call.ChannelId(channelId)
	}
	if hl != "" {
		call = call.Hl(hl)
	}
	call = call.MaxResults(maxResults)
	if mine != false {
		call = call.Mine(true)
	}
	if onBehalfOfContentOwner != "" {
		call = call.OnBehalfOfContentOwner(onBehalfOfContentOwner)
	}
	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
	if VideoUsecase.playlistId != "" {
		call = call.Id(VideoUsecase.playlistId)
	}
	response, err := call.Do()
	handleError(err, "")
	return response
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		logger.GetLogger().Error(message + ": " + err.Error())
	}
}
func (videoUsecase *VideoUsecase) GetVideos() ([]model.Video, error) {
	response := videoUsecase.channelsListMine([]string{"contentDetails"})
	var videos []model.Video

	for _, channel := range response.Items {
		playlistId := channel.ContentDetails.RelatedPlaylists.Uploads

		// Print the playlist ID for the list of uploaded videos.
		fmt.Printf("Videos in list %s\r\n", playlistId)

		nextPageToken := ""
		for {
			// Retrieve next set of items in the playlist.
			playlistResponse := videoUsecase.playlistItemsList([]string{"snippet"}, playlistId, nextPageToken)

			savePlaylistResponse(videoUsecase, playlistResponse, &videos)

			// Set the token to retrieve the next page of results
			// or exit the loop if all results have been retrieved.
			nextPageToken = playlistResponse.NextPageToken
			if nextPageToken == "" {
				break
			}
			fmt.Println()
		}
	}
	return videos, nil
}

func savePlaylistResponse(videoUsecase *VideoUsecase, playlistResponse *youtube.PlaylistItemListResponse, videos *[]model.Video) {
	for _, playlistItem := range playlistResponse.Items {
		title := playlistItem.Snippet.Title
		videoId := playlistItem.Snippet.ResourceId.VideoId
		fmt.Printf("%v, (%v)\r\n", title, videoId)

		// Insert video to database
		video := model.Video{
			YoutubeVideoID:         videoId,
			CreatedAt:              time.Now(),
			UpdatedAt:              time.Now(),
			CreatedBy:              1,
			UpdatedBy:              1,
			YoutubeTitle:           title,
			YoutubeDescription:     playlistItem.Snippet.Description,
			YoutubePlaylist:        playlistItem.Snippet.PlaylistId,
			YoutubeChannelID:       playlistItem.Snippet.ChannelId,
			YoutubeChannelUsername: playlistItem.Snippet.ChannelTitle,
		}

		if playlistItem.Status != nil {
			video.YoutubePrivacyStatus = playlistItem.Status.PrivacyStatus
		}
		result, err := videoUsecase.VideoRepository.GetVideoByVideoID(videoId)
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("Error while getting video by video id")
		}
		if result != nil {
			err = videoUsecase.VideoRepository.InsertVideo(result)
			if err != nil {
				logger.GetLogger().WithField("error", err).Error("Error while inserting video")
			}
			*videos = append(*videos, *result)
			continue
		}

		err = videoUsecase.VideoRepository.InsertVideo(&video)
		if err != nil {
			logger.GetLogger().WithField("error", err).Error("Error while inserting video")
		}
		*videos = append(*videos, video)
	}
}
