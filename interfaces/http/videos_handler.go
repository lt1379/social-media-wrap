package http

import (
	"github.com/gin-gonic/gin"
	"my-project/usecase"
)

type IVideoHandler interface {
	GetVideos(c *gin.Context)
}

type VideoHandler struct {
	VideoUsecase usecase.IVideoUsecase
}

func NewVideoHandler(videoUsecase usecase.IVideoUsecase) IVideoHandler {
	return &VideoHandler{VideoUsecase: videoUsecase}
}

func (videoHandler *VideoHandler) GetVideos(c *gin.Context) {
	res, err := videoHandler.VideoUsecase.GetVideos()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}
