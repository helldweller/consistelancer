package main

import (
    "context"
    "errors"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "golang.org/x/sync/errgroup"
    log "github.com/sirupsen/logrus"
    "package/main/db"
)

func main() {
    log.Info("Starting app")
    ctx, cancel := context.WithCancel(context.Background())
    group, groupCtx := errgroup.WithContext(ctx)
    msgChannel := make(chan db.Upstreams, 1)

    rdb := db.Initialize(groupCtx, config.RedisHost + ":" + config.RedisPort)
    defer rdb.Close(groupCtx)
    
    // goroutine to check for signals to gracefully finish all functions
    group.Go(func() error {
        signalChannel := make(chan os.Signal, 1)
        signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
        select {
        case sig := <-signalChannel:
            log.Error(fmt.Sprintf("Received signal: %s", sig))
            cancel()
        case <-groupCtx.Done():
            log.Error("Closing signal goroutine")
            return groupCtx.Err()
        }
        return nil
    })

    // goroutine to run k8sApiWatcher worker
    group.Go(func() error {
        interval := 3
        return k8sApiWatcher(interval, msgChannel, groupCtx)
    })

    // goroutine to run writeUpstreams worker
    group.Go(func() error {
        return writeUpstreams(msgChannel, rdb, groupCtx)
    })

    // goroutine to run readUpstreams worker
    group.Go(func() error {
        return readUpstreams(rdb, groupCtx)
    })
    
    // goroutine to run dbChecker worker
    group.Go(func() error {
        interval := 1
        return dbChecker(interval, groupCtx)
    })

    err := group.Wait()
    if err != nil {
        if errors.Is(err, context.Canceled) {
            log.Error("Context was canceled")
        } else {
            log.Error(fmt.Sprintf("Received error: %v\n", err))
        }
    } else {
        log.Error("Sucsessfull finished")
    }
}
