package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/models"
)

type (
	StatsService interface {
		UpdateGlobalStats(string, time.Duration)
		UpdateChannelsStats(string, time.Duration)
	}

	statsService struct{
		logger *zap.Logger
		pushStreamService PushStreamService
		redisClient redis.UniversalClient
		statsKeyPrefix string
	}
)

func (s *statsService) statsKey(suffix string) string {
	return fmt.Sprintf("%s_%s", s.statsKeyPrefix, suffix)
}

func (s *statsService) getStatsDetailed(ch chan map[string]interface{}) {
	data, err := s.pushStreamService.GetChannelsStatsDetailed()
	if err != nil {
		ch <- nil
		return
	}
	ch <- data
}

func (s *statsService) getStatsSummarized(ch chan map[string]interface{}) {
	data, err := s.pushStreamService.GetChannelsStatsSummarized()
	if err != nil {
		ch <- nil
		return
	}
	ch <- data
}

func (s *statsService) getStats() (map[string]interface{}, error) {
	chDetailed := make(chan map[string]interface{})
	chSummarized := make(chan map[string]interface{})

	go s.getStatsDetailed(chDetailed)
	go s.getStatsSummarized(chSummarized)

	detailed := <- chDetailed
	summarized := <- chSummarized

	if detailed == nil || summarized == nil {
		s.logger.Error("failed to get stats", zap.Any("detailed", detailed), zap.Any("summarized", summarized))
		return nil, errors.New("failed to get stats")
	}

	stats := map[string]interface{} {
		"detailed": detailed,
		"summarized": summarized,
	}

	return stats, nil
}

func (s *statsService) UpdateGlobalStats(keySuffix string, expiration time.Duration) {
	stats, err := s.getStats()
	if err != nil {
		return
	}

	key := s.statsKey(fmt.Sprintf("global_%s", keySuffix))
	value, err := json.Marshal(stats)
	err = s.redisClient.Set(key, value, expiration).Err()
	if err != nil {
		s.logger.Error("error saving global stats", zap.String("key", key), zap.Error(err))
		return
	}

	s.logger.Debug("did update global stats", zap.String("key", key))
}

func (s *statsService) UpdateChannelsStats(keySuffix string, expiration time.Duration) {
	data, err := s.pushStreamService.GetChannelsStatsDetailed()
	if err != nil {
		return
	}

	infos, ok := data["infos"]
	if !ok {
		s.logger.Error("data does not have property 'infos'")
		return
	}

	channels, ok := infos.([]interface{})
	if !ok {
		s.logger.Error("'infos' does not contain a list")
		return
	}

	pipeline := s.redisClient.Pipeline()
	for _, channel := range channels {
		var channelInfo models.ChannelInfo
		err := mapstructure.Decode(channel, &channelInfo)
		if err != nil {
			s.logger.Error("failed to map channel", zap.Error(err))
			continue
		}

		key := s.statsKey(fmt.Sprintf("channel_%s_host_%s", channelInfo.Channel, keySuffix))
		value, err := json.Marshal(channelInfo)
		pipeline.Set(key, value, expiration)
	}

	_, err = pipeline.Exec()
	if err != nil {
		s.logger.Error("failed to execute pipeline to update channels stats", zap.Error(err))
		return
	}
}

func NewStatsService(config *viper.Viper, logger *zap.Logger, redisClient redis.UniversalClient, pushStreamService PushStreamService) StatsService {
	statsKeyPrefix := config.GetString("redis.db.stats.prefix")

	return &statsService{
		logger: logger.Named("statsService"),
		pushStreamService: pushStreamService,
		redisClient: redisClient,
		statsKeyPrefix: statsKeyPrefix,
	}
}
