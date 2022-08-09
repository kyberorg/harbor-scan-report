package scan

import "github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/level"

type Report struct {
	Failed                  bool
	Scanner                 Scanner
	TopSeverity             level.FailLevel
	CriticalVulnerabilities []Vulnerability
	HighVulnerabilities     []Vulnerability
	MediumVulnerabilities   []Vulnerability
	LowVulnerabilities      []Vulnerability
	Counters                Counters
}

type Scanner struct {
	Name    string
	Vendor  string
	Version string
}

type Vulnerability struct {
	ID          string
	Url         string
	Package     string
	Version     string
	FixVersion  string
	Severity    level.FailLevel
	Description string
}

type Counters struct {
	Total    int
	Critical int
	High     int
	Medium   int
	Low      int
}
