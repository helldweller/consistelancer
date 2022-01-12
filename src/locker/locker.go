package locker

import (
    "context"
    "fmt"
    "time"
    // "math/rand"

    "github.com/go-redis/redis/v8"
    log "github.com/sirupsen/logrus"
    "github.com/bsm/redislock"
)


var (
    locker = &RedislockClient{}
    lockerKey = "lock"
)

type RedislockClient struct {
    rdbclient *redis.Client
    client *redislock.Client
    lock *redislock.Lock
    lockTtl time.Duration
}

func Initialize(ctx context.Context, ttl int, redisAddr string) *RedislockClient {
    rdbclient := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    if err := rdbclient.Ping(ctx).Err(); err != nil {
        log.Fatal("Redis is offline")
    }
    locker.lockTtl = time.Duration(ttl) * time.Second
    locker.rdbclient = rdbclient
    locker.client = redislock.New(rdbclient)
    return locker
}

func (locker *RedislockClient) Release(ctx context.Context) {
    locker.lock.Release(ctx)
}

func (locker *RedislockClient) Close() {
    locker.rdbclient.Close()
}

func (locker *RedislockClient) Obtain(ctx context.Context) error {
    lock, err := locker.client.Obtain(ctx, lockerKey, locker.lockTtl, nil)
    if err == nil {
        locker.lock = lock
    }
    return err
}

func (locker *RedislockClient) IsMasterOld(ctx context.Context) bool {
    result := false
    ttl, err := locker.lock.TTL(ctx)
    if (err == nil) || (ttl > 0) {
        if err := locker.lock.Refresh(ctx, locker.lockTtl, nil); err != nil {
            fmt.Println(err)
            return result
        }
        result = true
    }
    return result
}

func (locker *RedislockClient) IsMaster(ctx context.Context) (result bool) {
    result = false
    lock, err := locker.client.Obtain(ctx, lockerKey, locker.lockTtl, nil)
    if err != nil {
        return
    }
    locker.lock = lock
    ttl, err := locker.lock.TTL(ctx)
    if (err == nil) || (ttl > 0) {
        if err := locker.lock.Refresh(ctx, locker.lockTtl, nil); err != nil {
            log.Error(err)
            return
        }
        result = true
    }
    return
}