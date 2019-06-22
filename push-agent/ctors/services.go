package ctors

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/rafaeleyng/push-agent/push-agent/services"
)

func NewSubscriptionService(config *viper.Viper, logger *zap.Logger, pushStreamService services.PushStreamService, machineryServer *machinery.Server) services.SubscriptionService {
	return services.NewSubscriptionService(config, logger, pushStreamService, machineryServer)
}

func NewPushStreamService(config *viper.Viper, logger *zap.Logger) services.PushStreamService {
	return services.NewPushStreamService(config, logger)
}
