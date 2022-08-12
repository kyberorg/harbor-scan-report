package scan

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/severity"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/webutil"
	"io"
	"time"
)

const maxRetryAttempts = 3
const waitSeconds = 5

const scanReadyMarker = "Success"

var retryCounter = 0

func RunScan() *Report {
	var findResult *findImageOutput

	for {
		if retryCounter > maxRetryAttempts {
			err := fmt.Sprintf("Failed to get report after %d attempts\n", maxRetryAttempts+1)
			util.ExitOnError(errors.New(err))
			return nil
		} else if retryCounter > 0 {
			log.Info.Printf("Retry %d \n", retryCounter)
		}

		findResult = findImage()
		if findResult.Failed {
			util.ExitOnError(errors.New("search image request failed"))
			return nil
		}
		if !findResult.ImageFound {
			err := errors.New(fmt.Sprintf("Image '%s' is not found", config.Get().ImageInfo.Raw))
			util.ExitOnError(err)
			return nil
		}
		if findResult.ScanCompleted {
			log.Info.Println("Scan report is ready")
			break
		}

		log.Info.Printf("Scan report is not ready yet. Waiting for %d seconds", waitSeconds)
		time.Sleep(waitSeconds * time.Second)
		retryCounter++
	}
	return getScanReport(findResult)
}

func findImage() *findImageOutput {
	resp, err := webutil.DoFindRequest()
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	var goodResponse harbor.FindImageJson
	var errorResponse harbor.ErrorJson
	output := &findImageOutput{}

	//log.Debug.Println("Find image output: " + string(body))
	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &goodResponse)
		util.ExitOnError(err)

		output.Failed = false
		output.ImageFound = true

		if goodResponse == nil || len(goodResponse) == 0 {
			output.Failed = true
			log.Error.Println("Got malformed JSON in response. Raw response: " + string(body))
			return output
		}
		//log.Debug.Println("Find image resp: " + util.PrettyPrint(goodResponse))
		scanStatus := goodResponse[0].ScanOverview.VulnerabilityReport.ScanStatus
		output.ScanCompleted = util.IsStringPresent(scanStatus) && scanStatus == scanReadyMarker

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
			output.ImageFound = false
			break
		case 500:
			log.Error.Println("Harbor-Side error")
			output.Failed = true
		}
		err = json.Unmarshal(body, &errorResponse)
		if err == nil {
			log.Error.Println("Error searching image: " + errorResponse.Errors[0].Message)
		} else {
			log.Error.Println("Error searching image: " + string(body))
		}
		util.ExitOnError(err)
	}

	return output
}

func getScanReport(findResult *findImageOutput) *Report {
	scanResultsEndpoint := findResult.ScanResultsUrl
	//log.Debug.Println("Getting Scan Report from: " + scanResultsEndpoint)
	resp, err := webutil.DoScanReportRequest(scanResultsEndpoint)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	var goodResponse harbor.ScanResultsJson
	var errorResponse harbor.ErrorJson

	report := &Report{}

	//log.Debug.Println("Scan Report Raw body: " + string(body))
	if resp.StatusCode == 200 {
		report.Failed = false
		err = json.Unmarshal(body, &goodResponse)
		util.ExitOnError(err)
	} else {
		report.Failed = true
		err = json.Unmarshal(body, &errorResponse)
		if err == nil {
			log.Error.Println("Error getting scan report: " + errorResponse.Errors[0].Message)
		} else {
			log.Error.Println("Error getting scan report: " + string(body))
		}
		util.ExitOnError(err)
	}

	if !report.Failed {
		report = generateScanReport(goodResponse)
	}
	return report
}

func generateScanReport(json harbor.ScanResultsJson) *Report {
	report := &Report{
		Failed:                  false,
		GeneratedAt:             json.VulnerabilityReport.GeneratedAt,
		CriticalVulnerabilities: []Vulnerability{},
		HighVulnerabilities:     []Vulnerability{},
		MediumVulnerabilities:   []Vulnerability{},
		LowVulnerabilities:      []Vulnerability{},
	}
	report.Scanner = Scanner{
		Name:    json.VulnerabilityReport.Scanner.Name,
		Vendor:  json.VulnerabilityReport.Scanner.Vendor,
		Version: json.VulnerabilityReport.Scanner.Version,
	}
	if len(json.VulnerabilityReport.Vulnerabilities) > 0 {
		report.TopSeverity = severity.CreateFromString(json.VulnerabilityReport.Severity)
	} else {
		report.TopSeverity = severity.None
	}

	fixableVulnerabilityCounter := 0
	for _, v := range json.VulnerabilityReport.Vulnerabilities {
		currentSeverity := severity.CreateFromString(v.Severity)
		if currentSeverity.IsNotValid() {
			log.Warning.Printf("Skipping %s: wrong severity \n", v.ID)
		}
		vuln := Vulnerability{
			ID:          v.ID,
			Package:     v.Package,
			Version:     v.Version,
			FixVersion:  v.FixVersion,
			Severity:    currentSeverity,
			Description: v.Description,
		}
		if len(v.Links) > 0 {
			vuln.Url = v.Links[0]
		}
		if util.IsStringPresent(v.FixVersion) {
			fixableVulnerabilityCounter++
		}
		switch currentSeverity {
		case severity.Critical:
			report.CriticalVulnerabilities = append(report.CriticalVulnerabilities, vuln)
			break
		case severity.High:
			report.HighVulnerabilities = append(report.HighVulnerabilities, vuln)
			break
		case severity.Medium:
			report.MediumVulnerabilities = append(report.MediumVulnerabilities, vuln)
			break
		case severity.Low:
			report.LowVulnerabilities = append(report.LowVulnerabilities, vuln)
			break
		default:
			log.Warning.Printf("%s has unknown severity %s. Skipping.\n", vuln.ID, vuln.Severity)
		}
	}
	report.Counters = Counters{
		Total: len(report.CriticalVulnerabilities) + len(report.HighVulnerabilities) +
			len(report.MediumVulnerabilities) + len(report.LowVulnerabilities),
		Fixable:  fixableVulnerabilityCounter,
		Critical: len(report.CriticalVulnerabilities),
		High:     len(report.HighVulnerabilities),
		Medium:   len(report.MediumVulnerabilities),
		Low:      len(report.LowVulnerabilities),
	}
	return report
}
