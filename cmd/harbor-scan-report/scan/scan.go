package scan

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"io/ioutil"
	"net/http"
)

func RunScan() *Report {
	findResult := findImage()
	if findResult.Failed {
		util.ExitOnError(errors.New("search image request failed"))
	}
	if findResult.Found {
		return getScanReport(findResult)
	} else {
		err := errors.New(fmt.Sprintf("Image '%s' is not found", config.Get().ImageInfo.Raw))
		util.ExitOnError(err)
		return nil
	}
}

func findImage() *findImageOutput {
	findEndpoint := harbor.GetFindImageEndpoint()
	log.Debug.Println("Find Endpoint: " + findEndpoint)

	// Get request
	resp, err := http.Get(findEndpoint)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var goodResponse harbor.FindImageJson
	var errorResponse harbor.ErrorJson
	output := &findImageOutput{}

	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &goodResponse)
		util.ExitOnError(err)

		output.Failed = false
		output.Found = true

		if goodResponse[0].AdditionLinks.Vulnerabilities.Absolute {
			output.ScanResultsUrl = goodResponse[0].AdditionLinks.Vulnerabilities.Href
		} else {
			output.ScanResultsUrl = harbor.BuildHostPart() + goodResponse[0].AdditionLinks.Vulnerabilities.Href
		}
	} else {
		switch resp.StatusCode {
		case 400:
			log.Error.Println("Got Bad Request")
			output.Failed = true
			break
		case 401:
			log.Error.Println("Cannot request Harbor with given credentials: unauthorized. " +
				"Please check configuration again.")
			output.Failed = true
			break
		case 404:
			log.Error.Println("Image not found.")
			output.Failed = false
			output.Found = false
			break
		case 500:
			log.Error.Println("Harbor-Side error")
			output.Failed = true
		}
		err = json.Unmarshal(body, &errorResponse)
		util.ExitOnError(err)
		log.Error.Println("Error searching image: " + errorResponse.Errors[0].Message)
	}

	return output
}

func getScanReport(findResult *findImageOutput) *Report {
	scanResultsEndpoint := findResult.ScanResultsUrl
	log.Debug.Println("Getting Scan Report from: " + scanResultsEndpoint)
	//TODO implement
	return &Report{}
}
