package app

import (
    "context"
    "fmt"

    log "package/main/internal/logger"
    "package/main/internal/db"
)

func writeUpstreams(groupCtx context.Context, inputChannel chan db.Upstreams, rdb *db.RedisClient) error {
    log.Info("Starting writeUpstreams worker")
    for {
        select {
        case received := <-inputChannel:
            err := rdb.Publish(groupCtx, "upstreams", received)
            if err != nil {
                log.Error(fmt.Sprintf("Cant publish to redis: %s", err))
            }
        case <-groupCtx.Done():
            log.Error("Closing writeUpstreams goroutine")
            return groupCtx.Err()
        }
    }
}
