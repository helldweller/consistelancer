package app

import (
    "context"
    "fmt"
    "encoding/json"
    "reflect"

    log "package/main/internal/logger"
    "package/main/internal/db"
)

func readUpstreams(groupCtx context.Context, rdb *db.RedisClient) error {
    log.Info("Starting readUpstreams worker")
    for {
        select {
        case msg  := <-rdb.Channel:
            result := db.Upstreams{}
            if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
                log.Error(fmt.Sprintf("Cant Unmarshal received msg: %s", err))
            }
            if !reflect.DeepEqual(db.UpstreamList, result) {
                log.Info(fmt.Sprintf("Upstreams was updated: was %v, now %v", db.UpstreamList, result))
                db.UpstreamList = result
            }
            // item := db.UpstreamList.GetRandomItem()
            // log.Info(item)
        case <-groupCtx.Done():
            log.Error("Closing readUpstreams goroutine")
            return groupCtx.Err()
        }
    }
}
