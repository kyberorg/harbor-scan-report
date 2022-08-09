package webutil

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"net/http"
)

const (
	Accept                = "Accept"
	AcceptJson            = "application/json"
	ReportTypeHeaderName  = "X-Accept-Vulnerabilities"
	ReportTypeHeaderValue = "application/vnd.security.vulnerability.report; version=1.1"
)

func DoFindRequest() (*http.Response, error) {
	endpoint := harbor.GetFindImageEndpoint()
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(Accept, AcceptJson)
	if config.Get().Harbor.Credentials.Present {
		req.SetBasicAuth(config.Get().Harbor.Credentials.Robot, config.Get().Harbor.Credentials.Token)
	}
	return client.Do(req)
}

func DoScanReportRequest(endpoint string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(Accept, AcceptJson)
	req.Header.Add(ReportTypeHeaderName, ReportTypeHeaderValue)

	if config.Get().Harbor.Credentials.Present {
		req.SetBasicAuth(config.Get().Harbor.Credentials.Robot, config.Get().Harbor.Credentials.Token)
	}

	return client.Do(req)
}
