package scan

import (
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
)

func RunScan() *Report {
	findResult := findImage()
	if findResult.Found {
		return scanImage(findResult)
	} else {
		//TODO action
		err := errors.New(fmt.Sprintf("Image '%s' is not found", config.Get().ImageInfo.Raw))
		util.ExitOnError(err)
		return nil
	}
}

func findImage() *findImageOutput {
	//TODO implement
	return &findImageOutput{
		Found: true,
	}
}

func scanImage(findResult *findImageOutput) *Report {
	//TODO implement
	return &Report{}
}
