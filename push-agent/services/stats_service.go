package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
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
		channelStatsPrefix string
		globalStatsPrefix string
		logger            *zap.Logger
		pushStreamService PushStreamService
		redisClient       redis.UniversalClient
	}
)

func (s *statsService) globalStatsKey(suffix string) string {
	return fmt.Sprintf("%s_%s", s.globalStatsPrefix, suffix)
}

func (s *statsService) channelStatsKey(suffix string) string {
	return fmt.Sprintf("%s_%s", s.channelStatsPrefix, suffix)
}

func (s *statsService) getStatsDetailed(ch chan *models.GlobalStatsDetailed) {
	data, err := s.pushStreamService.GetGlobalStatsDetailed()
	if err != nil {
		ch <- nil
		return
	}
	ch <- data
}

func (s *statsService) getStatsSummarized(ch chan *models.GlobalStatsSummarized) {
	data, err := s.pushStreamService.GetGlobalStatsSummarized()
	if err != nil {
		ch <- nil
		return
	}
	ch <- data
}

func (s *statsService) getGlobalStats() (*models.GlobalStats, error) {
	chDetailed := make(chan *models.GlobalStatsDetailed)
	chSummarized := make(chan *models.GlobalStatsSummarized)

	go s.getStatsDetailed(chDetailed)
	go s.getStatsSummarized(chSummarized)

	detailed := <- chDetailed
	summarized := <- chSummarized

	if detailed == nil || summarized == nil {
		s.logger.Error("failed to get stats", zap.Any("detailed", detailed), zap.Any("summarized", summarized))
		return nil, errors.New("failed to get stats")
	}

	stats := models.GlobalStats{
		Detailed: detailed,
		Summarized: summarized,
	}

	return &stats, nil
}

func (s *statsService) UpdateGlobalStats(keySuffix string, expiration time.Duration) {
	stats, err := s.getGlobalStats()
	if err != nil {
		return
	}

	key := s.globalStatsKey(fmt.Sprintf("%s", keySuffix))
	value, err := json.Marshal(stats)
	if err != nil {
		s.logger.Error("error marshaling global stats", zap.Error(err))
		return
	}

	err = s.redisClient.Set(key, value, expiration).Err()
	if err != nil {
		s.logger.Error("error saving global stats", zap.String("key", key), zap.Error(err))
		return
	}

	s.logger.Debug("did update global stats", zap.String("key", key))
}

func (s *statsService) UpdateChannelsStats(keySuffix string, expiration time.Duration) {
	detailed, err := s.pushStreamService.GetGlobalStatsDetailed()
	if err != nil {
		return
	}

	// init pipeline
	pipeline := s.redisClient.Pipeline()
	defer func() {
		err := pipeline.Close()
		if err != nil {
			s.logger.Error("failed do close pipeline", zap.Error(err))
		}
	}()

	// fill pipeline commands
	for _, channelStats := range detailed.Infos {
		key := s.channelStatsKey(fmt.Sprintf("%s:%s", channelStats.Channel, keySuffix))
		value, err := json.Marshal(channelStats)
		if err != nil {
			s.logger.Error("error marshaling channel stats", zap.Error(err))
			return
		}
		pipeline.Set(key, value, expiration)
	}

	// exec pipeline
	_, err = pipeline.Exec()
	if err != nil {
		s.logger.Error("failed to execute pipeline to update channels stats", zap.Error(err))
		return
	}
}

func NewStatsService(config *viper.Viper, logger *zap.Logger, redisClient redis.UniversalClient, pushStreamService PushStreamService) StatsService {
	channelStatsPrefix := config.GetString("redis.db.stats_channel.prefix")
	globalStatsPrefix := config.GetString("redis.db.stats_global.prefix")

	return &statsService{
		channelStatsPrefix: channelStatsPrefix,
		globalStatsPrefix: globalStatsPrefix,
		logger:            logger.Named("statsService"),
		pushStreamService: pushStreamService,
		redisClient:       redisClient,
	}
}
