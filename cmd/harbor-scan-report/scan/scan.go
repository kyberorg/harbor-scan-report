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
	"io/ioutil"
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
	resp, err := webutil.DoFindRequest()
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
	log.Debug.Println("Getting Scan Report from: " + scanResultsEndpoint)
	resp, err := webutil.DoScanReportRequest(scanResultsEndpoint)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var goodResponse harbor.ScanResultsJson
	var errorResponse harbor.ErrorJson

	report := &Report{}

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
	report.TopSeverity = severity.CreateFromString(json.VulnerabilityReport.Severity)

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
