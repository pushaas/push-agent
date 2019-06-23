package services

import (
	"fmt"
	"net/http"

	"github.com/imroc/req"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/models"
)

type (
	PushStreamService interface {
		GetChannelsStatsDetailed() (map[string]interface{}, error)
		GetChannelsStatsSummarized() (map[string]interface{}, error)
		PublishMessage(*models.Message)
	}

	publicationService struct{
		config *viper.Viper
		reqClient *req.Req
		logger *zap.Logger
		publishEndpoint string
		statsEndpoint string
	}
)

func (s *publicationService) publishOnSingleChannel(channel string, content *string) {
	url := fmt.Sprintf("%s?id=%s", s.publishEndpoint, channel)

	header := make(http.Header)
	header.Set("Content-Type", "text/plain")
	_, err := s.reqClient.Post(url, *content, header)
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

func (s *publicationService) getStatsData(url string) (map[string]interface{}, error) {
	res, err := s.reqClient.Get(url)
	if err != nil {
		s.logger.Error("failed to get", zap.String("url", url), zap.Error(err))
		return nil, err
	}

	data := make(map[string]interface{})
	err = res.ToJSON(&data)
	if err != nil {
		s.logger.Error("failed to decode json", zap.String("url", url), zap.Error(err))
		return nil, err
	}

	return data, nil
}

func (s *publicationService) GetChannelsStatsDetailed() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s?id=ALL", s.statsEndpoint)
	return s.getStatsData(url)
}

func (s *publicationService) GetChannelsStatsSummarized() (map[string]interface{}, error) {
	url := s.statsEndpoint
	return s.getStatsData(url)
}

func NewPushStreamService(config *viper.Viper, logger *zap.Logger, reqClient *req.Req) PushStreamService {
	pushStreamAddr := config.GetString("push-stream.address")
	publishEndpoint := fmt.Sprintf("%s/pub", pushStreamAddr)
	statsEndpoint := fmt.Sprintf("%s/channels-stats", pushStreamAddr)

	return &publicationService{
		config: config,
		reqClient: reqClient,
		logger: logger.Named("pushStreamService"),
		publishEndpoint: publishEndpoint,
		statsEndpoint: statsEndpoint,
	}
}
