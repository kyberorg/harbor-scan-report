package scan

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/severity"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"time"
)

type Report struct {
	Failed                  bool
	GeneratedAt             time.Time
	Scanner                 Scanner
	Counters                Counters
	TopSeverity             severity.Severity
	AllVulnerabilities      []Vulnerability
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
	Score       float64
	Description string
}

func (v Vulnerability) HasFixVersion() bool {
	return util.IsStringPresent(v.FixVersion)
}

type Counters struct {
	Total    int
	Fixable  int
	Critical int
	High     int
	Medium   int
	Low      int
}
