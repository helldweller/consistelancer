package main

import (
    // "errors"
    "fmt"
    "os"

    log "github.com/sirupsen/logrus"
)

func init() {
    log.SetFormatter(&log.JSONFormatter{})
    log.SetOutput(os.Stdout)

    loglevel, err := log.ParseLevel(config.Loglevel)
    if err == nil {
        log.SetLevel(loglevel)
    } else {
        log.SetLevel(log.ErrorLevel)
        log.Error("Can`t parse LOG_LEVEL. Used default value: LOG_LEVEL=error")
    }
    log.Info(fmt.Sprintf("Used json logger and loglevel: %s", config.Loglevel))
}