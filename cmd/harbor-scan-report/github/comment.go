package github

import (
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
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
	b.WriteString(fmt.Sprintf("Result for image `%s` \n", config.Get().ImageInfo.Raw))
	b.WriteString(fmt.Sprintf("%d vulnerabilities found", report.Counters.Total))
	b.WriteString(fmt.Sprintf("[:no_entry:] %d critical [:fire:] %d high [:warning:] %d medium [:triangular_flag_on_post:] %d low",
		report.Counters.Critical,
		report.Counters.High,
		report.Counters.Medium,
		report.Counters.Low,
	))
	return b.String()
}
