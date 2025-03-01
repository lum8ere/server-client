package redis_manager

import (
	"context"
	"time"

	"backed-api-v2/libs/5_common/smart_context"

	"github.com/go-redis/redis/v8"
)

// RedisManager предоставляет базовые методы для работы с Redis.
type RedisManager struct {
	client *redis.Client
}

// NewRedisManager инициализирует подключение к Redis.
func NewRedisManager(sctx smart_context.ISmartContext) (*RedisManager, error) {
	addr := "localhost:6379"
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // по умолчанию без пароля
		DB:       0,  // используется базовый DB
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisManager{client: client}, nil
}

// SetValue записывает значение в Redis с указанным временем жизни.
func (rm *RedisManager) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rm.client.Set(ctx, key, value, expiration).Err()
}

// GetValue получает значение из Redis по ключу.
func (rm *RedisManager) GetValue(ctx context.Context, key string) (string, error) {
	return rm.client.Get(ctx, key).Result()
}

// Publish публикует сообщение в указанный канал.
func (rm *RedisManager) Publish(ctx context.Context, channel string, message interface{}) error {
	return rm.client.Publish(ctx, channel, message).Err()
}
