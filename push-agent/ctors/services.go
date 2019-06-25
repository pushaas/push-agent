package ctors

import (
	"github.com/go-redis/redis"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
)

func NewPushStreamService(config *viper.Viper, logger *zap.Logger, reqClient *req.Req) services.PushStreamService {
	return services.NewPushStreamService(config, logger, reqClient)
}

func NewSubscriptionService(logger *zap.Logger, pushStreamService services.PushStreamService) services.SubscriptionService {
	return services.NewSubscriptionService(logger, pushStreamService)
}

func NewStatsService(config *viper.Viper, logger *zap.Logger,  redisClient redis.UniversalClient, pushStreamService services.PushStreamService) services.StatsService {
	return services.NewStatsService(config, logger, redisClient, pushStreamService)
}
