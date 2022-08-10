package main

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/github"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
)

func main() {
	log.Debug.Println("Application Configuration: " + config.PrintConfig())

	//get scan results
	scanReport := scan.RunScan()
	log.Debug.Printf("Image '%s' Scan Report: \n %s", config.Get().ImageInfo.Raw, util.PrettyPrint(scanReport))

	log.Debug.Printf("Image '%s': %d vulnerabilities found (%d critical, %d high, %d medium, %d low)",
		config.Get().ImageInfo.Raw,
		scanReport.Counters.Total,
		scanReport.Counters.Critical,
		scanReport.Counters.High,
		scanReport.Counters.Medium,
		scanReport.Counters.Low,
	)

	//write comment
	if config.Get().Github.Enabled {
		github.WriteComment(scanReport)
	}

	if scanReport.TopSeverity.IsMoreCriticalThen(config.Get().MaxAllowedSeverity) {
		log.Error.Fatalf("Image has vulnerabilities that are more critical then allowed severity %s. "+
			"Check failed \n",
			config.Get().MaxAllowedSeverity.String())
	}
}
