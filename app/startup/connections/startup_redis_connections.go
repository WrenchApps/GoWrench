package connections

import (
	"context"
	"errors"
	"fmt"
	"wrench/app"
	"wrench/app/manifest/connection_settings"

	"github.com/redis/go-redis/v9"
)

var redisClients map[string]*redis.Client
var redisClusterClients map[string]*redis.ClusterClient

func loadConnectionsRedis(redisSettings []*connection_settings.RedisConnectionSettings) error {

	if len(redisSettings) > 0 && redisClients == nil {
		redisClients = make(map[string]*redis.Client)
	}

	if len(redisClusterClients) > 0 && redisClusterClients == nil {
		redisClusterClients = make(map[string]*redis.ClusterClient)
	}

	if len(redisSettings) > 0 {
		for _, setting := range redisSettings {

			if setting.IsCluster {
				redisClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
					Addrs:    setting.Addresses,
					Password: setting.Password,
				})

				if err := redisClusterClient.Ping(context.Background()).Err(); err != nil {
					app.LogError2(fmt.Sprintf("Error to connect to redis cluster | redis connection id %v", setting.Id), err)
					return err
				} else {
					app.LogInfo(fmt.Sprintf("Connected to redis cluster | redis connection id %v", setting.Id))
				}

			} else {
				redisClient := redis.NewClient(&redis.Options{
					Addr:     setting.Address,
					Password: setting.Password,
					DB:       setting.Db,
				})

				redisClients[setting.Id] = redisClient
				if err := redisClient.Ping(context.Background()).Err(); err != nil {
					app.LogError2(fmt.Sprintf("Error to connect to redis node | redis connection id %v", setting.Id), err)
					return err
				} else {
					app.LogInfo(fmt.Sprintf("Connected to redis node | redis connection id %v", setting.Id))
				}
			}
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

func GetRedisClusterConnection(redisConnectionId string) (*redis.ClusterClient, error) {
	if len(redisConnectionId) == 0 ||
		len(redisClusterClients) == 0 ||
		redisClusterClients[redisConnectionId] == nil {
		return nil, errors.New("redis without connection")
	}

	return redisClusterClients[redisConnectionId], nil
}
