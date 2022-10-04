package config

import "strings"

type SortCriteria string

const (
	Severity SortCriteria = "severity"
	Score    SortCriteria = "score"
)

func CreateSortCriteriaFromString(s string) SortCriteria {
	str := strings.ToLower(s)
	switch str {
	case "score":
		return Score
	default:
		return Severity
	}
}
