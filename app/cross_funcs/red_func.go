package cross_funcs

import (
	"wrench/app/startup/connections"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis"
	redisync "github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

var redsyncs map[string]*redsync.Redsync

func GetRedsyncInstance(redisConnectionId string) *redsync.Redsync {

	if len(redsyncs) == 0 {
		redsyncs = make(map[string]*redsync.Redsync)
	}

	rs := redsyncs[redisConnectionId]

	if rs == nil {
		var pool redis.Pool

		redisClient, _ := connections.GetRedisConnection(redisConnectionId)
		pool = redisync.NewPool(redisClient)

		rs = redsync.New(pool)
		redsyncs[redisConnectionId] = rs
	}

	return rs
}
