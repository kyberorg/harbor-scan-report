package github

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"os"
)

const StepSummaryEnvVar = "GITHUB_STEP_SUMMARY"

func WriteStepComment(scanReport *scan.Report) {
	//append to GITHUB_STEP_SUMMARY env
	report = scanReport
	message := createMessage()

	currentStepSummary := os.Getenv(StepSummaryEnvVar)
	err := os.Setenv(StepSummaryEnvVar, currentStepSummary+message)
	if err != nil {
		log.Warning.Printf("Failed to write Step Comment. Error: %s", err.Error())
	}
}
