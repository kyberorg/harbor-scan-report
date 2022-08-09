package level

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"strings"
)

type FailLevel int8

const (
	Critical  FailLevel = 4
	High      FailLevel = 3
	Medium    FailLevel = 2
	Low       FailLevel = 1
	None      FailLevel = 0
	Undefined FailLevel = -1
)

func CreateFromString(s string) FailLevel {
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

func CreateFromInt(i int) FailLevel {
	const MinLevel = 0
	const MaxLevel = 4
	if i < MinLevel || i > MaxLevel {
		return Undefined
	}
	switch i {
	case 0:
		return None
	case 1:
		return Low
	case 2:
		return Medium
	case 3:
		return High
	case 4:
		return Critical
	default:
		return Undefined
	}
}

func (level FailLevel) IsValid() bool {
	return level != Undefined
}

func (level FailLevel) IsNotValid() bool {
	return !level.IsValid()
}

func (level FailLevel) IsMoreCriticalThen(anotherLevel FailLevel) bool {
	return level < anotherLevel
}
