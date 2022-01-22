package logger

import (
    "fmt"
    "os"
    "github.com/sirupsen/logrus"
)

func Init(lvl string) {
    logrus.SetFormatter(&logrus.JSONFormatter{})
    logrus.SetOutput(os.Stdout)

    loglevel, err := logrus.ParseLevel(lvl)
    if err == nil {
        logrus.SetLevel(loglevel)
    } else {
        logrus.SetLevel(logrus.ErrorLevel)
        logrus.Error("Can`t parse LOG_LEVEL. Used default value: LOG_LEVEL=error")
    }

    logrus.Info(fmt.Sprintf("Used json logger and loglevel: %s", lvl))
}

func Info(msg string) {
    logrus.Info(msg)
}

func Error(msg string) {
    logrus.Error(msg)
}

func Fatal(msg string) {
    logrus.Fatal(msg)
}
