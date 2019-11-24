package cache

import (
	"encoding/json"
	"github.com/QSis/common/redis"
)

const (
	Short   = 60    // 1min
	Default = 600   // 10min
	Medium  = 3600  // 1h
	Long    = 86400 // 1day

	Prefix = "_c_"
)

type (
	Generator func(interface{}) error
)

func Indicate(key string, args ...interface{}) []interface{} {
	indicator := []interface{}{key}
	indicator = append(indicator, args...)
	return indicator
}

func Get(indicator interface{}, out interface{}, generator Generator, expire ...int) error {
	cacheContent := ""
	cacheName := ""
	nameBytes, err := json.Marshal(indicator)
	if err == nil {
		cacheName = Prefix + string(nameBytes)
		cacheContent, _ = redis.Get(cacheName)
	}
	if cacheContent != "" {
		return json.Unmarshal([]byte(cacheContent), out)
	}

	err = generator(out)
	if err == nil && cacheName != "" {
		cacheContent, err := json.Marshal(out)
		expireTime := Default
		if len(expire) > 0 {
			expireTime = expire[0]
		}
		if err == nil {
			redis.Set(cacheName, cacheContent, expireTime)
		}
	}
	return err
}

func Delete(indicator interface{}) error {
	cacheName := ""
	nameBytes, err := json.Marshal(indicator)
	if err == nil {
		cacheName = Prefix + string(nameBytes)
		redis.Del(cacheName)
	}
	return err
}
