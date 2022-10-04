package report

import (
	"fmt"
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
	if len(scanReport.Vulnerabilities) == 0 {
		return
	}
	fmt.Println("======= List of Vulnerabilities =======")
	for _, vuln := range scanReport.Vulnerabilities {
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
		fmt.Printf("%s %s. Score (CVSSv3): %.1f %s. Affected Package: %s %s. Fixed in: %s \n",
			cve, severity, score, severityFromCVSS, affectedPackage, vulnerableVersion, fixVersion)
	}
	fmt.Println("")
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
