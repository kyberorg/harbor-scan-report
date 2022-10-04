package main

import (
	"errors"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/github"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/image"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/report"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
)

func main() {
	log.Debug.Println("Application Configuration: " + config.PrintConfig())

	//find Image
	findResult := image.GetFinder().FindImage()
	if findResult.Failed() {
		if findResult.HasError() {
			util.ExitOnError(findResult.GetError())
		} else {
			util.ExitOnError(errors.New("failed to find image"))
		}
	}

	//get scan results
	scanStatus := scan.WaitForScanCompeted()
	scanReport := scan.GetScanReport(scanStatus)

	log.Debug.Printf("Image '%s' has %d vulnerabilities (%d critical, %d high, %d medium, %d low)\n",
		config.Get().ImageInfo.Raw,
		scanReport.Counters.Total,
		scanReport.Counters.Critical,
		scanReport.Counters.High,
		scanReport.Counters.Medium,
		scanReport.Counters.Low,
	)
	if scanReport.Counters.Total > 0 {
		log.Debug.Printf("%d/%d fixable\n", scanReport.Counters.Fixable, scanReport.Counters.Total)
	}

	//write comment
	if config.Get().Github.Enabled {
		github.WriteComment(scanReport)
	}

	report.WriteListOfVulnerabilities(scanReport)

	if scanReport.TopSeverity.IsMoreCriticalThen(config.Get().MaxAllowedSeverity) {
		var hasFixableVulnerabilities bool
		for _, vuln := range scanReport.AllVulnerabilities {
			if vuln.Severity.IsMoreCriticalThen(config.Get().MaxAllowedSeverity) {
				if vuln.HasFixVersion() {
					hasFixableVulnerabilities = true
					break
				}
			}
		}
		if hasFixableVulnerabilities {
			log.Error.Fatalf("Image has fixable vulnerabilities that are more critical "+
				"then allowed severity %s. "+
				"Check failed \n",
				config.Get().MaxAllowedSeverity.String())
		}
	}
}
