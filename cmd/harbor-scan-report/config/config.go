package config

import (
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/level"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"os"
	"strings"
)

const (
	DefaultProtocol = "https"
	DefaultTag      = "latest"
	DefaultLevel    = level.None
)

var (
	config  *appConfig
	err     error
	errText string
)

func Get() *appConfig {
	return config
}

func PrintConfig() string {
	return util.PrettyPrint(config)
}
func init() {
	config = &appConfig{
		Harbor: Harbor{
			Instance: HarborInstance{
				Host:       getHarborHost(),
				Protocol:   getHarborProto(),
				CustomPort: getHarborCustomPort(),
			},
			Credentials: HarborCredentials{
				Present: false,
				Robot:   getHarborRobot(),
				Token:   getHarborToken(),
			},
		},
		Github: Github{
			Enabled:    false,
			Token:      getGithubToken(),
			CommentUrl: getGithubCommentUrl(),
		},
		ImageInfo: ImageInfo{
			Raw:      getImage(),
			Project:  parseProject(),
			RepoName: parseRepo(),
			Tag:      parseTag(),
		},
		FailLevel: getFailLevel(),
	}
	updateCredentialsState()
	updateGitHubState()
}

func getHarborHost() string {
	host := os.Getenv("HARBOR_HOST")
	if util.IsStringEmpty(host) {
		err = errors.New("harbor host is undefined or empty")
		util.ExitOnError(err)
	}
	return host
}

func getHarborProto() string {
	protocol, isProtoPresent := os.LookupEnv("HARBOR_PROTO")
	if isProtoPresent {
		if util.IsStringEmpty(protocol) {
			protocol = DefaultProtocol
		} else {
			if protocol != "https" && protocol != "http" {
				errText = fmt.Sprintf("protocol %s is unsupported. Should be http or https", protocol)
				err = errors.New(errText)
				util.ExitOnError(err)
			}
		}
	} else {
		protocol = DefaultProtocol
	}
	return protocol
}

func getHarborCustomPort() string {
	port := os.Getenv("HARBOR_PORT")
	if util.IsStringPresent(port) && !util.IsPortValid(port) {
		errText = fmt.Sprintf("Harbor port %s is invalid. Please re-check your configuration.", port)
		err = errors.New(errText)
		util.ExitOnError(err)
	}
	return port
}

func getHarborRobot() string {
	robot := os.Getenv("HARBOR_ROBOT")
	if util.IsStringEmpty(robot) {
		log.Warning.Println("HARBOR_ROBOT is undefined. You can query public repositories only.")
	}
	return robot
}

func getHarborToken() string {
	token := os.Getenv("HARBOR_TOKEN")
	if util.IsStringEmpty(token) {
		log.Warning.Println("HARBOR_TOKEN is undefined. You can query public repositories only.")
	}
	return token
}

func updateCredentialsState() {
	config.Harbor.Credentials.Present = util.IsStringPresent(config.Harbor.Credentials.Robot) &&
		util.IsStringPresent(config.Harbor.Credentials.Token)
}

func getGithubToken() string {
	ghToken := os.Getenv("GITHUB_TOKEN")
	return ghToken
}

func getGithubCommentUrl() string {
	ghCommentUrl := os.Getenv("GITHUB_URL")
	return ghCommentUrl
}

func updateGitHubState() {
	config.Github.Enabled = util.IsStringPresent(config.Github.Token) && util.IsStringPresent(config.Github.CommentUrl)
}

func getImage() string {
	return os.Getenv("IMAGE")
}

func parseProject() string {
	image := parseImage()
	//separate project and repo
	projectRepoArr := strings.Split(image, "/")
	if len(projectRepoArr) < 2 {
		err = errors.New("image string is malformed. Format project/image:tag")
		util.ExitOnError(err)
	}
	return projectRepoArr[0]
}

func parseRepo() string {
	image := parseImage()
	//separate project and repo
	projectRepoArr := strings.Split(image, "/")
	if len(projectRepoArr) < 2 {
		err = errors.New("image string is malformed. Format project/image:tag")
		util.ExitOnError(err)
	}
	project := projectRepoArr[0]
	return strings.ReplaceAll(image, project+"/", "")
}

func parseTag() string {
	imageTagArray := splitImageString()
	if len(imageTagArray) == 1 {
		return DefaultTag
	} else {
		return imageTagArray[1]
	}
}

func getFailLevel() level.FailLevel {
	failLevelString := os.Getenv("FAIL_LEVEL")
	if util.IsStringEmpty(failLevelString) {
		return DefaultLevel
	}

	failLevel := level.CreateFromString(failLevelString)
	if failLevel.IsNotValid() {
		err = errors.New(fmt.Sprintf("Wrong fail level: %s \n", failLevelString))
		util.ExitOnError(err)
	}
	return failLevel
}

func parseImage() string {
	imageTagArray := splitImageString()
	return imageTagArray[0]
}

func splitImageString() []string {
	imageString := os.Getenv("IMAGE")
	if util.IsStringEmpty(imageString) {
		err = errors.New("image undefined or empty")
		util.ExitOnError(err)
	}
	//if image starts with registry name - cut it down
	registryName := getHarborHost()
	if strings.HasPrefix(imageString, registryName) {
		imageString = strings.Replace(imageString, registryName+"/", "", 1)
	}
	imageTagArray := strings.Split(imageString, ":")
	if len(imageTagArray) > 2 || len(imageTagArray) < 1 {
		err = errors.New("image string is malformed. Format project/image:tag")
		util.ExitOnError(err)
	}
	return imageTagArray
}
