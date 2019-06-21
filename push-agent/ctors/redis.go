package ctors

import (
	"github.com/go-redis/redis"
	"github.com/go-siris/siris/core/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewRedis(config *viper.Viper, logger *zap.Logger) (redis.UniversalClient, error) {
	log := logger.Named("redis")

	address := config.GetString("redis.address")
	password := config.GetString("redis.password")

	if address == "" {
		log.Error("address is required", zap.String("address", address))
		return nil, errors.New("redis address is required")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	return client, nil
}
