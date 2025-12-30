package connections

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"wrench/app"
	"wrench/app/converts"
	"wrench/app/manifest/connection_settings"

	"github.com/redis/go-redis/v9"
)

var redisClients map[string]redis.UniversalClient

func loadConnectionsRedis(redisSettings []*connection_settings.RedisConnectionSettings) error {

	if len(redisSettings) > 0 && redisClients == nil {
		redisClients = make(map[string]redis.UniversalClient)
	}

	if len(redisSettings) > 0 {
		for _, setting := range redisSettings {
			var tlsConfig *tls.Config = nil
			isTls, err := converts.ConvertStringToBool(setting.Tls)

			if err != nil {
				app.LogError2(fmt.Sprintf("Error to parse tls value | redis connection id %v", setting.Id), err)
				return err
			}

			if isTls {
				tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			}

			uClient := redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs:     setting.Addresses,
				DB:        setting.Db,
				Password:  setting.Password,
				TLSConfig: tlsConfig,
			})

			if err := uClient.Ping(context.Background()).Err(); err != nil {
				app.LogError2(fmt.Sprintf("Error to connect to redis | redis connection id %v", setting.Id), err)
				return err
			} else {
				app.LogInfo(fmt.Sprintf("Connected to redis | redis connection id %v", setting.Id))
			}

			redisClients[setting.Id] = uClient
		}
	}

	return nil
}

func GetRedisConnection(redisConnectionId string) (redis.UniversalClient, error) {
	if len(redisConnectionId) == 0 ||
		len(redisClients) == 0 ||
		redisClients[redisConnectionId] == nil {
		return nil, errors.New("redis without connection")
	}

	return redisClients[redisConnectionId], nil
}
