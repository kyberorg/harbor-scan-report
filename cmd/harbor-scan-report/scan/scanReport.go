package scan

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/severity"
)

type Report struct {
	Failed                  bool
	Scanner                 Scanner
	Counters                Counters
	TopSeverity             severity.Severity
	CriticalVulnerabilities []Vulnerability
	HighVulnerabilities     []Vulnerability
	MediumVulnerabilities   []Vulnerability
	LowVulnerabilities      []Vulnerability
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
	Severity    severity.Severity
	Description string
}

type Counters struct {
	Total    int
	Fixable  int
	Critical int
	High     int
	Medium   int
	Low      int
}
