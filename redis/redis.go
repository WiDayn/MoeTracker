package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

var RedisDb *redis.Client

func InitRedis() (err error) {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // 指定
		Password: "",
		DB:       2, // redis一共16个库，指定其中一个库即可
	})
	_, err = RedisDb.Ping().Result()
	return
}

func main() {
	err := InitRedis()
	if err != nil {
		fmt.Printf("connect redis failed! err : %v\n", err)
		return
	}
	fmt.Println("redis连接成功！")
}
