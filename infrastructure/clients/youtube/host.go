package youtube

type IYoutubeHost interface {
}

type YoutubeHost struct {
	id string
}

func NewYoutubeHost() IYoutubeHost {
	return &YoutubeHost{}
}
