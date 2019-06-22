package services

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/models"
)

type (
	PushStreamService interface {
		PublishMessage(*models.Message)
	}

	publicationService struct{
		config *viper.Viper
		logger *zap.Logger
		publishEndpoint string
	}
)

const contentType = "text/plain"


func (s *publicationService) publishOnSingleChannel(channel string, content *string) {
	url := fmt.Sprintf("%s?id=%s", s.publishEndpoint, channel)

	_, err := http.Post(url, contentType, bytes.NewBufferString(*content))
	if err != nil {
		s.logger.Error("error while posting on push-stream", zap.String("channel", channel), zap.String("content", *content))
		return
	}

	s.logger.Debug("published on push-stream", zap.String("channel", channel), zap.String("content", *content))
}

func (s *publicationService) PublishMessage(message *models.Message) {
	for _, channel := range message.Channels {
		go s.publishOnSingleChannel(channel, &message.Content)
	}
}

func NewPushStreamService(config *viper.Viper, logger *zap.Logger) PushStreamService {
	pushStreamAddr := config.GetString("push-stream.address")
	publishEndpoint := fmt.Sprintf("%s/pub", pushStreamAddr)

	return &publicationService{
		config: config,
		logger: logger.Named("pushStreamService"),
		publishEndpoint: publishEndpoint,
	}
}
