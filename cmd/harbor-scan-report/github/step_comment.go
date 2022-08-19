package github

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"os"
)

const StepSummaryEnvVar = "GITHUB_STEP_SUMMARY"

func WriteStepSummary(scanReport *scan.Report) {
	//append to GITHUB_STEP_SUMMARY env
	report = scanReport

	stepCommentFile := os.Getenv(StepSummaryEnvVar)
	if util.IsStringEmpty(stepCommentFile) {
		log.Warning.Printf("Skipping writing Step Summary. %s envVar is empty", StepSummaryEnvVar)
		return
	}

	message := "## It works!"

	err := os.WriteFile(stepCommentFile, []byte(message), 0700)
	if err != nil {
		log.Warning.Printf("Failed to write step summary. Got I/O error: %s \n", err.Error())
	} else {
		log.Info.Println("Step Summary is written")
	}
}
