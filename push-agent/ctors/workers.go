package ctors

import (
	"strings"

	"github.com/Pallinder/sillyname-go"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
	"github.com/rafaeleyng/push-agent/push-agent/workers"
)

func NewStatsWorker(config *viper.Viper, logger *zap.Logger, statsService services.StatsService) workers.StatsWorker {
	agentName := strings.ReplaceAll(strings.ToLower(sillyname.GenerateStupidName()), " ", "-")
	return workers.NewStatsWorker(config, logger, agentName, statsService)
}

func NewSubscriptionWorker(config *viper.Viper, logger *zap.Logger, redisClient redis.UniversalClient, subscriptionService services.SubscriptionService) workers.SubscriptionWorker {
	return workers.NewSubscriptionWorker(config, logger, redisClient, subscriptionService)
}
