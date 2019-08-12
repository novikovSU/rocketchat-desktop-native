package utils

import "log"

func AssertErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func AssertErrMsg(err error, msg string) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
