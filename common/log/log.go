package log

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func Info(origin, function, msg interface{}) {
	log.Info(fmt.Sprintf("%s - %s => msg: %s", origin, function, msg))
}

func Warning(origin, function, msg interface{}) {
	log.Warning(fmt.Sprintf("%s - %s => warning: %s", origin, function, msg))
}

func Error(origin, function, msg interface{}) {
	log.Error(fmt.Sprintf("%s - %s => error: %s", origin, function, msg))
}

func Debug(origin, function, msg interface{}) {
	log.Debug(fmt.Sprintf("%s - %s => debug: %s", origin, function, msg))
}

func Fatal(origin, function, msg interface{}) {
	log.Fatal(fmt.Sprintf("%s - %s => fatal: %s", origin, function, msg))
}
