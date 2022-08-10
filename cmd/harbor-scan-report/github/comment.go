package github

import (
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/webutil"
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
	return fmt.Sprintf("Image '%s': %d vulnerabilities found (%d critical, %d high, %d medium, %d low)",
		config.Get().ImageInfo.Raw,
		report.Counters.Total,
		report.Counters.Critical,
		report.Counters.High,
		report.Counters.Medium,
		report.Counters.Low,
	)
}
