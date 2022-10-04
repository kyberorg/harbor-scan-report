package report

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
)

const (
	CriticalSeverity = "CRITICAL"
	HighSeverity     = "HIGH"
	MediumSeverity   = "MEDIUM"
	LowSeverity      = "LOW"
	InvalidSeverity  = "UNDEFINED"
)

func WriteListOfVulnerabilities(scanReport *scan.Report) {
	if len(scanReport.AllVulnerabilities) == 0 {
		return
	}
	log.Info.Println("List of Vulnerabilities")
	for _, vuln := range scanReport.AllVulnerabilities {
		cve := vuln.ID
		severity := vuln.Severity
		score := vuln.Score
		severityFromCVSS := severityFromScore(score)
		affectedPackage := vuln.Package
		vulnerableVersion := vuln.Version

		var fixVersion string
		if vuln.HasFixVersion() {
			fixVersion = vuln.FixVersion
		} else {
			fixVersion = "UNFIXABLE"
		}
		log.Info.Printf("%s %s Score (CVSS): %s %s. Affected Package: %s %s. Fixed in: %s \n",
			cve, severity, score, severityFromCVSS, affectedPackage, vulnerableVersion, fixVersion)
	}
}

func severityFromScore(score float64) string {
	if score > 10 || score < 0 {
		return InvalidSeverity
	} else if score >= 9 {
		return CriticalSeverity
	} else if score >= 7 {
		return HighSeverity
	} else if score >= 4 {
		return MediumSeverity
	} else {
		return LowSeverity
	}
}
