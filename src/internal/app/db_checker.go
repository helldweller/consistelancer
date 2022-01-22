package app

import (
    "context"
    "time"

    log "package/main/internal/logger"
    "package/main/internal/config"
    "package/main/internal/db"
)

func dbChecker(groupCtx context.Context, interval int, conf *config.Config) error {
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    log.Info("Starting dbChecker")
    for {
        select {
        case <-ticker.C:
            if err := db.Check(groupCtx); err != nil {
                log.Error(err.Error())
                db.Initialize(groupCtx, conf.RedisHost + ":" + conf.RedisPort)
            }
        case <-groupCtx.Done():
            log.Error("Closing dbChecker goroutine")
            return groupCtx.Err()
        }
    }
}
