package harbor

import (
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"strings"
)

const ApiTwoZero = "api/v2.0"

func GetFindImageEndpoint() string {
	project := config.Get().ImageInfo.Project
	repo := config.Get().ImageInfo.RepoName
	tag := config.Get().ImageInfo.Tag

	if strings.Contains(repo, "/") {
		repo = strings.ReplaceAll(repo, "/", "%252%")
	}

	return fmt.Sprintf("%s/%s/projects/%s/repositories/%s/artifacts?q=tags%%3D%s",
		BuildHostPart(), ApiTwoZero, project, repo, tag)
}

func BuildHostPart() string {
	proto := config.Get().Harbor.Instance.Protocol
	host := config.Get().Harbor.Instance.Host
	customPort := config.Get().Harbor.Instance.CustomPort

	url := proto + "://" + host
	if util.IsStringPresent(customPort) {
		url = url + ":" + customPort
	}
	return url
}
