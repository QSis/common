package redis

import (
	"github.com/QSis/common/config"
	"gopkg.in/redis.v5"
	"time"
)

var (
	Client    *redis.Client
	RedisConf *config.Config
)

func InitRedisWithConfig() {
	cfg := RedisConf
	if cfg == nil {
		cfg, _ := config.Config.Get("redis")
	}
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.UString("addr"),
		Password: cfg.UString("auth", ""),
	})
}

func Get(key string) (string, error) {
	return Client.Get(key).Result()
}

func Set(key string, value interface{}, expire ...int) error {
	expireTime := 0
	if len(expire) > 0 {
		expireTime = expire[0]
	}
	return Client.Set(key, value, time.Duration(expireTime)*time.Second).Err()
}

func SetNX(key string, value interface{}, expire ...int) (bool, error) {
	expireTime := 0
	if len(expire) > 0 {
		expireTime = expire[0]
	}
	return Client.SetNX(key, value, time.Duration(expireTime)*time.Second).Result()
}

func Expire(key string, expire int) (bool, error) {
	return Client.Expire(key, time.Duration(expire)*time.Second).Result()
}

func LPush(key string, value ...interface{}) (int, error) {
	len, err := Client.LPush(key, value...).Result()
	return int(len), err
}

func RPush(key string, value ...interface{}) (int, error) {
	len, err := Client.RPush(key, value...).Result()
	return int(len), err
}

func LPop(key string) (string, error) {
	return Client.LPop(key).Result()
}

func RPop(key string) (string, error) {
	return Client.RPop(key).Result()
}

func LRange(key string, start int64, stop int64) ([]string, error) {
	return Client.LRange(key, start, stop).Result()
}

func Pull(key string) ([]string, error) {
	pipe := Client.TxPipeline()
	cmd := pipe.LRange(key, 0, -1)
	pipe.Del(key)
	_, err := pipe.Exec()
	return cmd.Val(), err
}

func HMSet(key string, value map[string]string) error {
	return Client.HMSet(key, value).Err()
}

func HSet(key string, field string, value interface{}) error {
	return Client.HSet(key, field, value).Err()
}

func HGet(key string, field string) (string, error) {
	return Client.HGet(key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return Client.HGetAll(key).Result()
}

func HDel(key string, field ...string) (int64, error) {
	return Client.HDel(key, field...).Result()
}

func Del(key ...string) (int64, error) {
	return Client.Del(key...).Result()
}

func Rename(key string, newKey string) error {
	return Client.Rename(key, newKey).Err()
}
