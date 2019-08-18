package utils

import (
	log "github.com/chaykin/log4go"
)

func CreateLogger(category string) *log.Filter {
	return log.LOGGER(category)
}

func init() {
	log.LoadConfiguration("./log4go.json")
}
