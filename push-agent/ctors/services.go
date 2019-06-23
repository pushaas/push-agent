package ctors

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/go-redis/redis"
	"github.com/imroc/req"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
)

func NewPushStreamService(config *viper.Viper, logger *zap.Logger, reqClient *req.Req) services.PushStreamService {
	return services.NewPushStreamService(config, logger, reqClient)
}

func NewSubscriptionService(config *viper.Viper, logger *zap.Logger, pushStreamService services.PushStreamService, machineryServer *machinery.Server) services.SubscriptionService {
	return services.NewSubscriptionService(config, logger, pushStreamService, machineryServer)
}

func NewStatsService(config *viper.Viper, logger *zap.Logger,  redisClient redis.UniversalClient, pushStreamService services.PushStreamService) services.StatsService {
	return services.NewStatsService(config, logger, redisClient, pushStreamService)
}
