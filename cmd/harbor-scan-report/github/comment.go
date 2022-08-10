package github

import (
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/webutil"
	"strings"
)

func WriteComment(report *scan.Report) {
	message := createMessage(report)
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

func createMessage(report *scan.Report) string {
	var b strings.Builder

	b.WriteString("## Harbor Scan Image Report \n")
	b.WriteString(fmt.Sprintf("Result for image [%s](%s) \n", config.Get().ImageInfo.Raw, harbor.UiUrl()))
	b.WriteString(fmt.Sprintf("Total %d vulnerabilities found (%d fixable) \n",
		report.Counters.Total, report.Counters.Fixable))
	b.WriteString(fmt.Sprintf("Total %d vulnerabilities "+
		"("+
		"[:no_entry:](## \"critical\") %d critical "+
		"[:fire:](## \"high\") %d high "+
		"[:warning:](## \"medium\") %d medium "+
		"[:triangular_flag_on_post:](## \"low\") %d low"+
		")",
		report.Counters.Total,
		report.Counters.Critical,
		report.Counters.High,
		report.Counters.Medium,
		report.Counters.Low,
	))
	return b.String()
}
