package cache

import (
	"context"
	"github.com/go-redis/redis/v9"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strconv"
)

var (
	RedisClient *redis.Client

	RedisDb    string
	RedisAddr    string
	RedisPw    string
	RedisDbName    string

	)

func init()  {
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		logging.Fatalf("ini load redis file failed, err:%s \n", err)
	}

	LoadRedis(file)
	ConnectRedis()
}

func ConnectRedis()  {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		Password: RedisPw,  // 未设置密码 注释掉就好了
		DB: int(db),
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logging.Fatalf("redis ping failed, err:%s \n", err)
	}

	RedisClient = client
	logging.Println("connect redis success")
}

func LoadRedis(file *ini.File)  {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}