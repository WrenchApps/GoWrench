package connections

import (
	"errors"
	"wrench/app/manifest/connection_settings"

	"github.com/redis/go-redis/v9"
)

var redisClients map[string]*redis.Client

func loadConnectionsRedis(redisSettings []*connection_settings.RedisConnectionSettings) error {

	if len(redisSettings) > 0 && redisClients == nil {
		redisClients = make(map[string]*redis.Client)
	}

	if len(redisSettings) > 0 {
		for _, setting := range redisSettings {

			redisClient := redis.NewClient(&redis.Options{
				Addr:     setting.Address,
				Password: setting.Password,
				DB:       setting.Db,
			})

			redisClients[setting.Id] = redisClient
		}
	}

	return nil
}

func GetRedisConnection(redisConnectionId string) (*redis.Client, error) {
	if len(redisConnectionId) == 0 ||
		len(redisClients) == 0 ||
		redisClients[redisConnectionId] == nil {
		return nil, errors.New("redis without connection")
	}

	return redisClients[redisConnectionId], nil
}
