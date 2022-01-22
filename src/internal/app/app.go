package app

import (
    "context"
    "errors"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "golang.org/x/sync/errgroup"

	log "package/main/internal/logger"
    "package/main/internal/db"
	"package/main/internal/config"
)

func Run() {
	conf, err := config.Init()
    if err != nil && conf == nil {
        log.Fatal(err.Error())
    }
    if err != nil && conf != nil {
        log.Error(err.Error())
    }
    
    log.Init(conf.Loglevel)
    log.Info("Starting app")

    ctx, cancel := context.WithCancel(context.Background())
    group, groupCtx := errgroup.WithContext(ctx)
    upstreamChannel := make(chan db.Upstreams, 1)
	defer close(upstreamChannel)

    rdb := db.Initialize(groupCtx, conf.RedisHost + ":" + conf.RedisPort)
    defer rdb.Close(groupCtx)
    
    // goroutine to check for signals to gracefully finish all functions
    group.Go(func() error {
        signalChannel := make(chan os.Signal, 1)
        defer close(signalChannel)
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
        return k8sApiWatcher(groupCtx, interval, upstreamChannel, conf)
    })

    // goroutine to run writeUpstreams worker
    group.Go(func() error {
        return writeUpstreams(groupCtx, upstreamChannel, rdb)
    })

    // goroutine to run readUpstreams worker
    group.Go(func() error {
        return readUpstreams(groupCtx, rdb)
    })
    
    // goroutine to run dbChecker worker
    group.Go(func() error {
        interval := 1
        return dbChecker(groupCtx, interval, conf)
    })

    err = group.Wait()
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
