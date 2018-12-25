package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func InitLogging() {
	var logLevel log.Level

	defaultLogLevel := log.InfoLevel
	log.SetLevel(defaultLogLevel)

	customLoglevel := os.Getenv("LOG_LEVEL")
	if customLoglevel != "" {
		var err error
		logLevel, err = log.ParseLevel(customLoglevel)
		if err != nil {
			log.Warnf("Couldn't parse '%s'", customLoglevel)
		} else {
			log.SetLevel(logLevel)
		}
	}
}
