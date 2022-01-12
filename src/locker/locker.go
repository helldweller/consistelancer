package locker

import (
    "context"

	goredislib "github.com/go-redis/redis/v8"
    log "github.com/sirupsen/logrus"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
    
)

var (
    locker = &RedislockClient{}
    mutexname = "lock"
)

type RedislockClient struct {
    rdbclient *goredislib.Client
    mutex *redsync.Mutex
    // lockTtl time.Duration
}

func Initialize(ctx context.Context, redisAddr string) *RedislockClient {
	locker.rdbclient = goredislib.NewClient(&goredislib.Options{
		Addr: redisAddr,
	})
    if err := locker.rdbclient.Ping(ctx).Err(); err != nil {
        log.Fatal("Redis is offline")
    }
	pool := goredis.NewPool(locker.rdbclient)
	rs := redsync.New(pool)
	locker.mutex = rs.NewMutex(mutexname)
	return locker
}

func (locker *RedislockClient) Close() {
    locker.mutex.Unlock()
	locker.rdbclient.Close()
}

func (locker *RedislockClient) IsMaster() bool {
	if ok, err := locker.mutex.Extend(); ok && err == nil {
		return true
	}
	if err := locker.mutex.Lock(); err == nil {
		return true
	}
	return false
}
