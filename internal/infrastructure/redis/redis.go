package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/config"
	"github.com/spf13/cast"
)

var (
	once sync.Once
	sets *RedisSets
)

const (
	CHAT_INSTANCE_NAME = "chat" //list redis
)

type Redis struct {
	Host string
	Port string
	Pwd  string
	Db   string
	Name string
}

type RedisSets struct {
	redis map[string]*redis.Client
	l     sync.RWMutex
}

func (r *RedisSets) getRedis(key ...string) *redis.Client {
	r.l.RLock()
	defer r.l.RUnlock()
	name := "default"
	if len(key) > 0 {
		name = key[0]
	}
	if client, ok := r.redis[name]; ok {
		return client
	}
	return nil
}

func NewRedis() *RedisSets {
	once.Do(func() {
		redisSets := map[string]*redis.Client{}
		RedisConfig := config.GetGlobalConf().Redis
		for _, r := range RedisConfig {
			client := redis.NewClient(&redis.Options{
				Addr:       fmt.Sprintf("%s:%d", r.Host, r.Port),
				Password:   r.Pwd,
				DB:         cast.ToInt(r.Db),
				MaxRetries: 3, //重试次数
			})
			_, err := client.Ping(context.Background()).Result()
			if err != nil {
				panic("redis初始化失败:" + err.Error())
			}
			redisSets[r.Name] = client
		}
		sets = &RedisSets{
			redis: redisSets,
		}
	})
	return sets
}

func getRedis(keys ...string) *redis.Client {
	rds := NewRedis()
	key := ""
	if len(keys) > 0 {
		key = keys[0]
	}
	return rds.getRedis(key)
}
func (r *RedisSets) GetChatRedis() *redis.Client {
	return getRedis(CHAT_INSTANCE_NAME)
}

func GetChatRedis() *redis.Client {
	return getRedis(CHAT_INSTANCE_NAME)
}

func baseLock(ctx context.Context, key string, lockTime time.Duration, redisCli *redis.Client) (success bool, unlock func(), err error) {

	maxTryNum := 2000
	tryTime := 3 * time.Millisecond
	value := time.Now().UnixMicro()
	val := strconv.FormatInt(value, 10)
	stm := false
	for i := 0; i < maxTryNum; i++ {
		if stm {
			break
		}
		stm, err = redisCli.SetNX(ctx, key, val, lockTime).Result()
		if err != nil {
			panic(err)
		}
		time.Sleep(tryTime)
	}
	if stm == false {
		return false, nil, nil
	}
	return true, func() {
		c := context.Background()
		result, err := redisCli.Get(c, key).Result()
		if err != nil {
			return
		}
		if result != val {
			return
		}
		redisCli.Del(c, key)
	}, nil
}
