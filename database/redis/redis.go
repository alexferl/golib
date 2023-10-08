package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func New() (*redis.Client, error) {
	var opt *redis.Options
	var err error

	addr := viper.GetString(RedisAddr)
	if addr == "" {
		addr = DefaultConfig.Addr
	}

	url := viper.GetString(RedisURL)
	if url == "" {
		opt = &redis.Options{
			Addr:       viper.GetString(RedisAddr),
			Username:   viper.GetString(RedisUsername),
			Password:   viper.GetString(RedisPassword),
			ClientName: viper.GetString(RedisClientName),
			DB:         viper.GetInt(RedisDB),
		}
	} else {
		opt, err = redis.ParseURL(url)
		if err != nil {
			return nil, err
		}
	}

	client := redis.NewClient(opt)

	return client, nil
}
