package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	return fmt.Sprintf("%s:%s", s.globalStatsPrefix, suffix)
}

func (s *statsService) channelStatsKey(suffix string) string {
	return fmt.Sprintf("%s:%s", s.channelStatsPrefix, suffix)
}

func (s *statsService) UpdateGlobalStats(keySuffix string, expiration time.Duration) {
	// get data
	stats, err := s.pushStreamService.GetGlobalStatsSummarized()
	if err != nil {
		return
	}

	// prepare data
	// TODO remove if not needed
	//stats.Updated = time.Now().UTC()
	value, err := json.Marshal(stats)
	if err != nil {
		s.logger.Error("error marshaling global stats", zap.Error(err))
		return
	}

	// save on redis
	key := s.globalStatsKey(keySuffix)
	err = s.redisClient.Set(key, value, expiration).Err()
	if err != nil {
		s.logger.Error("error saving global stats", zap.String("key", key), zap.Error(err))
		return
	}

	s.logger.Debug("did update global stats", zap.String("key", key))
}

func (s *statsService) UpdateChannelsStats(keySuffix string, expiration time.Duration) {
	// get data
	stats, err := s.pushStreamService.GetGlobalStatsDetailed()
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
	for _, channelStats := range stats.Infos {
		// TODO remove if not needed
		//channelStats.Updated = time.Now().UTC()
		channelStats.Hostname = stats.Hostname

		value, err := json.Marshal(channelStats)
		if err != nil {
			s.logger.Error("error marshaling channel stats", zap.Error(err))
			return
		}

		key := s.channelStatsKey(fmt.Sprintf("%s:%s", channelStats.Channel, keySuffix))
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
