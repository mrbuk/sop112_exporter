package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// InitLogging initilizes the log subsystem with the default
// log level or one provided via environment variable LOG_LEVEL
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
