package ctors

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func getRedisUrl(config *viper.Viper) string {
	return config.GetString("redis.url")
}

func getSentinelMasterName(config *viper.Viper) string {
	return config.GetString("redis.sentinel.master_name")
}

func convertToUniversalOptions(options *redis.Options) *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:              []string{options.Addr},
		Password:           options.Password,
		DB:                 options.DB,
		MaxRetries:         options.MaxRetries,
		OnConnect:          options.OnConnect,
		MinRetryBackoff:    options.MinRetryBackoff,
		MaxRetryBackoff:    options.MaxRetryBackoff,
		DialTimeout:        options.DialTimeout,
		ReadTimeout:        options.ReadTimeout,
		WriteTimeout:       options.WriteTimeout,
		PoolSize:           options.PoolSize,
		MinIdleConns:       options.MinIdleConns,
		MaxConnAge:         options.MaxConnAge,
		PoolTimeout:        options.PoolTimeout,
		IdleTimeout:        options.IdleTimeout,
		IdleCheckFrequency: options.IdleCheckFrequency,
		TLSConfig:          options.TLSConfig,
	}
}

func getRedisOptions(config *viper.Viper) (*redis.UniversalOptions, error) {
	redisUrl := getRedisUrl(config)
	options, err := redis.ParseURL(redisUrl)

	if err != nil {
		urlOptions, err := url.Parse(redisUrl)

		if err != nil {
			return nil, errors.New("failed to parse redis URL")
		}

		password, _ := urlOptions.User.Password()
		addresses := strings.Split(urlOptions.Host, ",")

		return &redis.UniversalOptions{
			MasterName: getSentinelMasterName(config),
			Addrs:      addresses,
			Password:   password,
		}, nil
	}

	if options.Addr == "" {
		return nil, errors.New("redis address is required")
	}

	universalOptions := convertToUniversalOptions(options)

	return universalOptions, nil
}

func NewRedisClient(config *viper.Viper, logger *zap.Logger) (redis.UniversalClient, error) {
	log := logger.Named("redisClient")

	options, err := getRedisOptions(config)
	if err != nil {
		log.Error("failed to init redis options", zap.Error(err))
		return nil, err
	}

	log.Info("initializing redis with address", zap.String("addr", options.Addrs[0]))
	client := redis.NewUniversalClient(options)

	return client, nil
}
