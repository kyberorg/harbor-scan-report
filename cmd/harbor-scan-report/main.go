package main

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/github"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
)

func main() {
	log.Info.Println("Starting the application...")
	log.Info.Println("Logging system initialized")
	log.Debug.Println(config.PrintConfig())

	//get scan results
	scan.RunScan()
	//write comment
	if config.Get().Github.Enabled {
		github.WriteComment()
	}
}
