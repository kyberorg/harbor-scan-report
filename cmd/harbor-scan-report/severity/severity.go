package severity

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"strings"
)

type Severity int8

const (
	Critical  Severity = 4
	High      Severity = 3
	Medium    Severity = 2
	Low       Severity = 1
	None      Severity = 0
	Undefined Severity = -1
)

func CreateFromString(s string) Severity {
	if util.IsStringEmpty(s) {
		return Undefined
	}
	str := strings.ToLower(s)
	switch str {
	case "critical":
		return Critical
	case "high":
		return High
	case "medium":
		return Medium
	case "low":
		return Low
	case "none":
		return None
	default:
		return Undefined
	}
}

func (s Severity) IsValid() bool {
	return s != Undefined
}

func (s Severity) IsNotValid() bool {
	return !s.IsValid()
}

func (s Severity) IsMoreCriticalThen(anotherSeverity Severity) bool {
	return s > anotherSeverity
}

func (s Severity) String() string {
	switch s {
	case None:
		return "None"
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return "Undefined"
	}
}
