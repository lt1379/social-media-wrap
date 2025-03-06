package youtube_client

import (
	"fmt"
	"google.golang.org/api/youtube/v3"
)

type ITestYoutubeClient interface {
	ChannelsListByUsername(part []string, forUsername string)
}

type TestYoutubeClient struct {
	YoutubeClient *youtube.Service
}

func NewTestYoutubeClient(youtubeClient *youtube.Service) ITestYoutubeClient {
	return &TestYoutubeClient{YoutubeClient: youtubeClient}
}

func (testYoutubeClient *TestYoutubeClient) ChannelsListByUsername(part []string, forUsername string) {
	call := testYoutubeClient.YoutubeClient.Channels.List(part)
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}
