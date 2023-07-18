package caching

import (
	"cais/pkg/logger"
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

type RedisEnvars struct {
	Addr string
	Pass string
	Db   string
}

// NewRedisClient подключенает redis, возвращает клиент
func NewRedisClient(ctx context.Context, log *logger.Logger, re RedisEnvars) (rdCl *redis.Client, err error) {
	addr := os.Getenv(re.Addr)

	passwd := os.Getenv(re.Pass)

	db, err := strconv.Atoi(os.Getenv(re.Db))
	if err != nil {
		log.Errorf("env get redis db error: %v", err)
		return nil, err
	}

	cl := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	pong, err := cl.Ping(ctx).Result()
	if err != nil {
		log.Errorf("test connection (ping-pong) error: %v", err)
		return nil, err
	}
	log.Infof("test ping - %s(result) redis connection", pong)

	return cl, err
}
