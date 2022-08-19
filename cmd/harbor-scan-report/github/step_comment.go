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

	currentStepSummary := os.Getenv(StepSummaryEnvVar)
	err := os.Setenv(StepSummaryEnvVar, currentStepSummary+"## It works!")
	if err != nil {
		log.Warning.Printf("Failed to write Step Comment. Error: %s", err.Error())
	}

	log.Debug.Printf("Current value: %s /n", os.Getenv(StepSummaryEnvVar))
}
