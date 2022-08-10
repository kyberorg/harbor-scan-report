package github

import (
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/severity"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/webutil"
	"strings"
)

var report *scan.Report

func WriteComment(scanReport *scan.Report) {
	report = scanReport
	message := createMessage()
	resp, err := webutil.DoGitHubCommentRequest(message)
	if err != nil {
		log.Warning.Printf("Failed to create GitHub Comment")
	}
	if resp.StatusCode == 201 {
		log.Info.Println("GitHub comment created")
	} else {
		log.Warning.Printf("Failed to create GitHub comment. Status: %d \n", resp.StatusCode)
	}
}

func createMessage() string {
	var b strings.Builder

	b.WriteString("## Harbor Image Vulnerability Report \n")
	b.WriteString(fmt.Sprintf("Results for image [%s](%s) \n", config.Get().ImageInfo.Raw, harbor.UiUrl()))
	b.WriteString(fmt.Sprintf("Total %d vulnerabilities found ",
		report.Counters.Total))
	if report.Counters.Total > 0 {
		b.WriteString(fmt.Sprintf("- %d fixable ", report.Counters.Fixable))
	} else {
		b.WriteString(s2e(severity.None))
	}
	b.WriteString(fmt.Sprintf("\n"))
	if report.Counters.Total > 0 {
		b.WriteString(fmt.Sprintf("[%s](## \"total vulnerabilities\") %d vulnerabilities "+
			"("+
			"[%s](## \"critical\") %d critical "+
			"[%s](## \"high\") %d high "+
			"[%s](## \"medium\") %d medium "+
			"[%s](## \"low\") %d low"+
			")",
			topSeverityEmoji(), report.Counters.Total,
			s2e(severity.Critical), report.Counters.Critical,
			s2e(severity.High), report.Counters.High,
			s2e(severity.Medium), report.Counters.Medium,
			s2e(severity.Low), report.Counters.Low,
		))
	}
	return b.String()
}

func s2e(s severity.Severity) string {
	switch s {
	case severity.Critical:
		return ":no_entry:"
	case severity.High:
		return ":fire:"
	case severity.Medium:
		return ":warning:"
	case severity.Low:
		return ":triangular_flag_on_post:"
	case severity.None:
		return ":heavy_check_mark:"
	default:
		return ":interrobang:"
	}
}

func topSeverityEmoji() string {
	return s2e(report.TopSeverity)
}
