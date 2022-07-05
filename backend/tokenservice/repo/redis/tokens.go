package tokens

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type redisTokenRepository struct {
	*redis.Client
}

func NewRedisTokenRepository(host, port, password string) *redisTokenRepository {
	return &redisTokenRepository{
		Client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       0,
		}),
	}
}

func (rdb *redisTokenRepository) SaveToken(userID, tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	return rdb.Set(key, 1, expiration).Err()
}

func (rdb *redisTokenRepository) IsTokenValid(userID, tokenID string) (bool, error) {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	res, err := rdb.Get(key).Result()
	if err != nil {
		return false, err
	}
	if res != "1" {
		return false, nil
	}

	return true, nil
}

func (rdb *redisTokenRepository) InvalidateToken(userID, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	res, err := rdb.Get(key).Result()
	if err != nil {
		return err
	}
	if res != "1" {
		return errors.New("token already invalid")
	}
	if err := rdb.Do("set", key, 0, "keepttl").Err(); err != nil {
		return err
	}

	return nil
}
