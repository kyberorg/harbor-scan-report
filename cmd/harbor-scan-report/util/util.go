package util

import (
	"encoding/json"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"strconv"
	"strings"
	"time"
)

func IsStringEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsStringPresent(s string) bool {
	return !IsStringEmpty(s)
}

func ExitOnError(err error) {
	if err != nil {
		log.Error.Fatalln(err.Error())
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

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func PrettyDate(t time.Time) string {
	return t.Format("2.1.2006 15:04:05-0700")
}
