package services

import (
	"fmt"
	"net/http"

	"github.com/imroc/req"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/models"
)

type (
	PushStreamService interface {
		GetGlobalStatsDetailed() (*models.GlobalStatsDetailed, error)
		GetGlobalStatsSummarized() (*models.GlobalStatsSummarized, error)
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

func (s *publicationService) GetGlobalStatsDetailed() (*models.GlobalStatsDetailed, error) {
	url := fmt.Sprintf("%s?id=ALL", s.statsEndpoint)

	data, err := s.getStatsData(url)
	if err != nil {
		return nil, err
	}

	var stats models.GlobalStatsDetailed
	err = mapstructure.Decode(data, &stats)
	if err != nil {
		s.logger.Error("failed to decode detailed data", zap.String("url", url), zap.Error(err))
		return nil, err
	}

	return &stats, nil
}

func (s *publicationService) GetGlobalStatsSummarized() (*models.GlobalStatsSummarized, error) {
	url := s.statsEndpoint

	data, err := s.getStatsData(url)
	if err != nil {
		return nil, err
	}

	var stats models.GlobalStatsSummarized
	err = mapstructure.Decode(data, &stats)
	if err != nil {
		s.logger.Error("failed to decode sumarized data", zap.String("url", url), zap.Error(err))
		return nil, err
	}

	return &stats, nil
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
