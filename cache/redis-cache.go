package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
)

type RedisCache struct {
	host string
	db   int
	exp  time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) *RedisCache {
	return &RedisCache{host, db, exp}
}

func (c *RedisCache) GetClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.host,
		Password: "", // no password set
		DB:       c.db, // use default DB
	})
}

func (c *RedisCache) Set(key string, value interface{}) {
	client := c.GetClient()
	json,err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	client.Set(key, json, c.exp)
}

func (c *RedisCache) Get(key string) (interface{}, bool) {
	client := c.GetClient()
	val, err := client.Get(key).Result()
	if err != nil {
		return nil, false
	}
	var result interface{}
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		panic(err)
	}
	return result, true
}