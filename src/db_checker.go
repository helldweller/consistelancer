package main

import (
    "context"
    // "fmt"
    "time"

    log "github.com/sirupsen/logrus"
    "package/main/db"
)

func dbChecker(interval int, groupCtx context.Context) error {
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    log.Info("Starting dbChecker")
    for {
        select {
        case <-ticker.C:
            if err := db.Check(groupCtx); err != nil {
                log.Error(err.Error())
                db.Initialize(groupCtx, config.RedisHost + ":" + config.RedisPort)
            }
        case <-groupCtx.Done():
            log.Error("Closing dbChecker goroutine")
            return groupCtx.Err()
        }
    }
}
