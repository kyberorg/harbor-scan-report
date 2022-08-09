package config

import (
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/severity"
)

type appConfig struct {
	Harbor             Harbor
	Github             Github
	ImageInfo          ImageInfo
	MaxAllowedSeverity severity.Severity
}

type Harbor struct {
	Instance    HarborInstance
	Credentials HarborCredentials
}

type HarborInstance struct {
	Host       string
	Protocol   string
	CustomPort string
}

type HarborCredentials struct {
	Present bool
	Robot   string
	Token   string
}

type Github struct {
	Enabled    bool
	Token      string
	CommentUrl string
}

type ImageInfo struct {
	Raw      string
	Project  string
	RepoName string
	Tag      string
}
