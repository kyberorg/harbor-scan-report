package webutil

import (
	"bytes"
	"encoding/json"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"net/http"
	"strconv"
)

const (
	Accept        = "Accept"
	Authorization = "Authorization"
	GitHubJson    = "application/vnd.github+json"
	Json          = "application/json"

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
	req.Header.Add(Accept, Json)
	req.Header.Add(ReportTypeHeaderName, ReportTypeHeaderValue)
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
	req.Header.Add(Accept, Json)
	req.Header.Add(ReportTypeHeaderName, ReportTypeHeaderValue)

	if config.Get().Harbor.Credentials.Present {
		req.SetBasicAuth(config.Get().Harbor.Credentials.Robot, config.Get().Harbor.Credentials.Token)
	}

	return client.Do(req)
}

func DoGitHubCommentCreateRequest(comment string) (*http.Response, error) {
	ghComment := githubComment{Body: comment}
	body, err := json.Marshal(ghComment)
	if err != nil {
		return nil, err
	}
	endpoint := config.Get().Github.CommentUrl
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(Accept, GitHubJson)
	req.Header.Add(Authorization, "token "+config.Get().Github.Token)

	return client.Do(req)
}

func DoGitHubCommentSearchRequest() (*http.Response, error) {
	endpoint := config.Get().Github.CommentUrl
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(Accept, GitHubJson)
	req.Header.Add(Authorization, "token "+config.Get().Github.Token)

	return client.Do(req)
}

func DoGitHubCommentUpdateRequest(commentId int, comment string) (*http.Response, error) {
	ghComment := githubComment{Body: comment}
	body, err := json.Marshal(ghComment)
	if err != nil {
		return nil, err
	}
	endpoint := config.Get().Github.CommentUrl + "/" + strconv.Itoa(commentId)
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(Accept, GitHubJson)
	req.Header.Add(Authorization, "token "+config.Get().Github.Token)

	return client.Do(req)
}

type githubComment struct {
	Body string `json:"body"`
}
