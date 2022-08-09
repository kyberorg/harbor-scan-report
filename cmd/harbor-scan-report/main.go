package main

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/github"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
)

func main() {
	config.PrintConfig()

	//get scan results
	scan.RunScan()
	//write comment
	if config.Get().Github.Enabled {
		github.WriteComment()
	}
}
