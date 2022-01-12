package main

import (
    "context"
    "fmt"
    "encoding/json"
    "reflect"

    log "github.com/sirupsen/logrus"
    "package/main/db"
)

func readUpstreams(rdb *db.RedisClient, groupCtx context.Context) error {
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
        case <-groupCtx.Done():
            log.Error("Closing readUpstreams goroutine")
            return groupCtx.Err()
        }
    }
}
