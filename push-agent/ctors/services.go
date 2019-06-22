package ctors

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
)

func NewSubscriptionService(config *viper.Viper, logger *zap.Logger, redisClient redis.UniversalClient, pushStreamService services.PushStreamService) services.SubscriptionService {
	return services.NewSubscriptionService(config, logger, redisClient, pushStreamService)
}

func NewPushStreamService(config *viper.Viper, logger *zap.Logger) services.PushStreamService {
	return services.NewPushStreamService(config, logger)
}
