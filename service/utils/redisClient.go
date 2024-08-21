package utils

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	Rdb = redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_IP"),
		Password: GetEnv("REDIS_PWD"),
		DB:       0, // use default DB
	})
)

func InitRedis() bool {
	PrintObj(GetEnv("REDIS_IP"), "InitRedis")
	ticker := time.NewTicker(1 * time.Second)
	count := 0
	success := false

	for range ticker.C {
		initRedisClient()

		if err := PingRedis(); err == nil {
			success = true
		}

		PrintObj("InitRedis retry:"+ToString(count), "")

		if count == 60 || success {
			return success
		}

		count++
	}

	return true
}

func initRedisClient() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     GetEnv("REDIS_IP"),
		Password: GetEnv("REDIS_PWD"),
		DB:       0, // use default DB
	})
}

func PingRedis() error {
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	return nil
}

func SetRedis(key string, value string, hours int16) error {
	time := time.Duration(hours) * time.Hour
	err := Rdb.Set(ctx, key, value, time).Err()
	PrintObj([]string{key, fixStringLen(value)}, "SetRedis")

	if err != nil {
		PrintObj(err, "SetRedis err")
	}

	return err
}

func GetRedis(key string) string {
	value, err := Rdb.Get(ctx, key).Result()
	PrintObj([]string{key, fixStringLen(value)}, "GetRedis")

	if err != nil {
		PrintObj(err, "GetRedis err")
		return ""
	}

	return value
}

func GetRedisKeys(prefix string) []string {
	var cursor uint64
	res := []string{}

	for {
		var keys []string
		var err error
		keys, cursor, err = Rdb.Scan(ctx, cursor, prefix, 0).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			res = append(res, key)
		}

		if cursor == 0 { // no more keys
			break
		}
	}

	PrintObj(res, "GetRedisKeys")

	return res
}

func HasRedis(key string) bool {
	val, err := Rdb.Exists(ctx, key).Result()
	PrintObj([]interface{}{key, val}, "HasRedis")

	if err != nil || val == 0 {
		PrintObj("redis not found", "")
		return false
	}

	return true
}

func DeleteRedis(key string) {
	PrintObj(key, "DeleteRedis")
	_, err := Rdb.Del(ctx, key).Result()

	if err != nil {
		PrintObj(err, "DeleteRedis err")
	}
}
