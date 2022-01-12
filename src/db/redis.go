package db

import (
    "encoding/json"
    "time"
    // "fmt"
    "context"

    "github.com/go-redis/redis/v8"
    log "github.com/sirupsen/logrus"
)

var (
    rdb = &RedisClient{}
)

type RedisClient struct {
    client *redis.Client
    // pubsub *redis.PubSub
    Channel <-chan *redis.Message
}

func Initialize(ctx context.Context, addr string) *RedisClient {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    if err := client.Ping(ctx).Err(); err != nil {
        log.Error("Unable to connect to redis " + err.Error())
    }
    rdb.client = client
    pubsub := rdb.client.Subscribe(ctx, "upstreams")
    rdb.Channel = pubsub.Channel()
    return rdb
}

func Close() {
    // func (c *PubSub) Unsubscribe(ctx context.Context, channels ...string) error
    rdb.client.Close()
}

func Check(ctx context.Context) error {
    return rdb.client.Ping(ctx).Err()
}

func (rdb *RedisClient) GetKey(ctx context.Context, key string, src interface{}) error {
    val, err := rdb.client.Get(ctx, key).Result()
    if err == redis.Nil || err != nil {
        return err
    }
    err = json.Unmarshal([]byte(val), &src)
    if err != nil {
        return err
    }
    return nil
}

func (rdb *RedisClient) SetKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    err := rdb.client.Set(ctx, key, value, expiration).Err() // cacheEntry
    if err != nil {
        return err
    }
    return nil
}

func (rdb *RedisClient) Publish(ctx context.Context, key string, value interface{}) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    err = rdb.client.Publish(ctx, key, data).Err()
    if err != nil {
        return err
    }
    return nil
}
