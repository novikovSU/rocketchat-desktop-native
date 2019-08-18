package utils

import (
	"os"

	log "github.com/chaykin/log4go"
)

var (
	logger *log.Filter
)

func init() {
	logger = CreateLogger("main")
}

func AssertErr(err error) {
	if err != nil {
		logger.Critical(err)
		os.Exit(1)
	}
}

func AssertErrMsg(err error, msg string) {
	if err != nil {
		logger.Critical(msg, err)
		os.Exit(1)
	}
}

func Safe(obj interface{}, err error) {
	if err != nil {
		logger.Critical("Could not perform operation on object %s. Cause: %s\n", obj, err)
		os.Exit(1)
	}
}
