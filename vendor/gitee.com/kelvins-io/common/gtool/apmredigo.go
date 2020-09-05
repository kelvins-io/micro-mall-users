package gtool

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"go.elastic.co/apm/module/apmredigo"
)

func WithRedigoContext(redisPool *redis.Pool, ctx context.Context) redis.Conn {
	return apmredigo.Wrap(redisPool.Get()).WithContext(ctx)
}
