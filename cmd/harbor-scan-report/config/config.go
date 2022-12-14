package config

import (
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/comment"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/severity"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultProtocol      = "https"
	DefaultTag           = "latest"
	DefaultLevel         = severity.None
	DefaultCommentTitle  = "Docker Image Vulnerability Report"
	DefaultTimeout       = 120
	DefaultCheckInterval = 5
	DefaultSortCriteria  = Severity
)

var (
	config  *appConfig
	timer   *Timer
	err     error
	errText string
)

func Get() *appConfig {
	return config
}

func IsTimeoutSet() bool {
	return config.Timing.Timeout != 0
}

func GetTimer() *Timer {
	return timer
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
			Digest:   parseDigest(),
		},
		MaxAllowedSeverity: getMaxAllowedSeverity(),
		Comment: Comment{
			Title: getCommentTitle(),
			Mode:  getCommentMode(),
		},
		Timing: Timing{
			Timeout:       getTimeout(),
			CheckInterval: getCheckInterval(),
		},
		Report: Report{
			SortBy:          getSortCriteria(),
			ShowFixableOnly: getShowFixableOnly(),
		},
	}
	updateCredentialsState()
	updateGitHubState()
	createTimer()
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
	return strings.TrimSpace(ghCommentUrl)
}

func updateGitHubState() {
	config.Github.Enabled = util.IsStringPresent(config.Github.Token) && util.IsStringPresent(config.Github.CommentUrl)
}

func createTimer() {
	timer = &Timer{
		SecondsLeft: config.Timing.Timeout,
	}
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
		tagPart := imageTagArray[1]
		if strings.Contains(tagPart, "@") {
			tagArray := strings.Split(tagPart, "@")
			if len(tagArray) > 0 {
				return tagArray[0]
			} else {
				err = errors.New("image string is malformed. Format project/image:tag[@digest]")
				util.ExitOnError(err)
				return ""
			}

		}
		return imageTagArray[1]
	}
}

func parseDigest() string {
	imageString := getImageString()
	imageShaArray := strings.Split(imageString, "@")

	if len(imageShaArray) > 2 || len(imageShaArray) < 1 {
		err = errors.New("image string is malformed. Format project/image:tag[@digest]")
		util.ExitOnError(err)
		return ""
	} else if len(imageShaArray) == 2 {
		return imageShaArray[1]
	} else {
		//no sha in image
		digestParam := os.Getenv("DIGEST")
		if util.IsStringPresent(digestParam) && strings.HasPrefix(digestParam, "sha256:") {
			return digestParam
		} else {
			return ""
		}
	}
}

func getMaxAllowedSeverity() severity.Severity {
	failLevelString := os.Getenv("MAX_ALLOWED_SEVERITY")
	if util.IsStringEmpty(failLevelString) {
		return DefaultLevel
	}

	failLevel := severity.CreateFromString(failLevelString)
	if failLevel.IsNotValid() {
		err = errors.New(fmt.Sprintf("Wrong fail level: %s \n", failLevelString))
		util.ExitOnError(err)
	}
	return failLevel
}

func getCommentTitle() string {
	customTitle := os.Getenv("COMMENT_TITLE")
	if util.IsStringEmpty(customTitle) {
		customTitle = DefaultCommentTitle
	}
	return customTitle
}

func getCommentMode() comment.Mode {
	commentMode := os.Getenv("COMMENT_MODE")
	err, mode := comment.CreateCommentMode(commentMode)
	util.ExitOnError(err)
	return mode
}

func getTimeout() int {
	timeoutString := os.Getenv("TIMEOUT")
	if util.IsStringEmpty(timeoutString) {
		return DefaultTimeout
	}
	timeout, err := strconv.Atoi(timeoutString)
	if err != nil {
		util.ExitOnError(errors.New("timeout should be number"))
	}

	if timeout < 0 {
		err = errors.New("timeout should be positive number of seconds")
		util.ExitOnError(err)
	}
	return timeout
}

func getCheckInterval() int {
	intervalString := os.Getenv("CHECK_INTERVAL")
	if util.IsStringEmpty(intervalString) {
		return DefaultCheckInterval
	}
	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		util.ExitOnError(errors.New("check interval should be number"))
	}
	if interval < 0 {
		err = errors.New("interval should be positive number of seconds")
		util.ExitOnError(err)
	}
	return interval
}

func getSortCriteria() SortCriteria {
	sortCriteriaString := os.Getenv("REPORT_SORT_BY")
	if util.IsStringEmpty(sortCriteriaString) {
		return DefaultSortCriteria
	}
	return CreateSortCriteriaFromString(sortCriteriaString)
}

func getShowFixableOnly() bool {
	fixableOnlyString := os.Getenv("REPORT_ONLY_FIXABLE")
	if util.IsStringEmpty(fixableOnlyString) {
		return false
	}
	fixableOnly, err := strconv.ParseBool(fixableOnlyString)
	if err != nil {
		err = errors.New("parameter 'report-only-fixable' should have only boolean value (true/false)")
		util.ExitOnError(err)
	}
	return fixableOnly
}

func parseImage() string {
	imageTagArray := splitImageString()
	return imageTagArray[0]
}

func splitImageString() []string {
	imageString := getImageString()
	//if image starts with registry name - cut it down
	registryName := getHarborHost()
	if strings.HasPrefix(imageString, registryName) {
		imageString = strings.Replace(imageString, registryName+"/", "", 1)
	}
	imageTagArray := strings.Split(imageString, ":")
	if len(imageTagArray) > 3 || len(imageTagArray) < 1 {
		err = errors.New("image string is malformed. Format project/image:tag[@digest]")
		util.ExitOnError(err)
	}
	return imageTagArray
}

func getImageString() string {
	imageString := os.Getenv("IMAGE")
	if util.IsStringEmpty(imageString) {
		err = errors.New("image undefined or empty")
		util.ExitOnError(err)
	}
	return imageString
}
