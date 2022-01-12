package main

import (
    "context"
    "fmt"

    log "github.com/sirupsen/logrus"
    "package/main/db"
)

func writeUpstreams(inputChannel chan db.Upstreams, rdb *db.RedisClient, groupCtx context.Context) error {
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
