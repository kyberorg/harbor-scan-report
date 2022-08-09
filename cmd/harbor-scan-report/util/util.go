package util

import (
	"log"
	"strconv"
	"strings"
)

func IsStringEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsStringPresent(s string) bool {
	return !IsStringEmpty(s)
}

func ExitOnError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func IsPortValid(portString string) bool {
	port, err := strconv.Atoi(portString)
	if err != nil {
		return false
	} else {
		const MinPort = 1
		const MaxPort = 65535
		return port >= MinPort && port <= MaxPort
	}
}
