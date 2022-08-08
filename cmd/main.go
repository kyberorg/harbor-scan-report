package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type HarborInstance struct {
	Host     string
	Protocol string
	Port     string
}

type HarborAccess struct {
	Robot string
	Token string
}

type Github struct {
	Enabled    bool
	Token      string
	CommentUrl string
}

type Image struct {
	Project  string
	RepoName string
	Tag      string
}
type FailLevel int8

const (
	Critical FailLevel = 4
	High     FailLevel = 3
	Medium   FailLevel = 2
	Low      FailLevel = 1
	None     FailLevel = 0
)

const (
	DefaultProtocol = "https"
	DefaultTag      = "latest"
)

var (
	err     error
	errText string
)

func main() {
	//get args
	println("HH: " + os.Getenv("HARBOR_HOST"))
	println("GH URL: " + os.Getenv("GITHUB_URL"))
	println("Fail level: " + os.Getenv("FAIL_LEVEL"))

	harborInstance, err := createHarborInstance()
	exitOnError(err)

	harborAccess, err := createHarborAccess()
	exitOnError(err)

	githubObject := createGitHubObject()

	image, err := createImageObject()
	exitOnError(err)

	//get scan results
	runScan(harborInstance, harborAccess, image)
	//write comment
	if githubObject.Enabled {
		writeComment(githubObject)
	}
}

func createHarborInstance() (*HarborInstance, error) {
	//Harbor Host
	host := os.Getenv("HARBOR_HOST")
	if isStringEmpty(host) {
		err = errors.New("harbor host is undefined or empty")
		return nil, err
	}
	//Protocol
	protocol, isProtoPresent := os.LookupEnv("HARBOR_PROTO")
	if isProtoPresent {
		if isStringEmpty(protocol) {
			protocol = DefaultProtocol
		} else {
			if protocol != "https" && protocol != "http" {
				errText = fmt.Sprintf("protocol %s is unsupported. Should be http or https", protocol)
				err = errors.New(errText)
				return nil, err
			}
		}
	} else {
		protocol = DefaultProtocol
	}
	port := os.Getenv("HARBOR_PORT")
	if !isPortValid(port) {
		errText = fmt.Sprintf("port %s is invalid. Please re-check your configuration.", port)
		err = errors.New(errText)
		return nil, err
	}
	return &HarborInstance{
		Host:     host,
		Protocol: protocol,
		Port:     port,
	}, nil
}

func createHarborAccess() (*HarborAccess, error) {
	robot := os.Getenv("HARBOR_ROBOT")
	if isStringEmpty(robot) {
		return nil, errors.New("please specify Harbor Robot (or Username) to access Harbor with")
	}
	token := os.Getenv("HARBOR_TOKEN")
	if isStringEmpty(token) {
		return nil, errors.New("please specify Harbor Token (or Password) to access Harbor with")
	}
	return &HarborAccess{
		Robot: robot,
		Token: token,
	}, nil
}

func createGitHubObject() *Github {
	ghToken := os.Getenv("GITHUB_TOKEN")
	ghCommentUrl := os.Getenv("GITHUB_URL")

	ghEnabled := isStringPresent(ghToken) && isStringPresent(ghCommentUrl)

	return &Github{
		Enabled:    ghEnabled,
		Token:      ghToken,
		CommentUrl: ghCommentUrl,
	}
}

func createImageObject() (*Image, error) {
	var project string
	var repo string
	var tag string
	var image string

	imageString := os.Getenv("IMAGE")
	if isStringEmpty(imageString) {
		return nil, errors.New("image undefined or empty")
	}
	//separate image and tag
	imageTagArray := strings.Split(imageString, ":")
	if len(imageTagArray) > 2 || len(imageTagArray) < 1 {
		return nil, errors.New("image string is malformed. Format project/image:tag")
	} else if len(imageTagArray) == 2 {
		image = imageTagArray[0]
		tag = imageTagArray[1]
	} else if len(imageTagArray) == 1 {
		image = imageTagArray[0]
		tag = DefaultTag
	}
	//separate project and repo
	projectRepoArr := strings.Split(image, "/")
	if len(projectRepoArr) < 2 {
		return nil, errors.New("image string is malformed. Format project/image:tag")
	} else {
		project = projectRepoArr[0]
		repo = strings.ReplaceAll(image, project+"/", "")
	}
	return &Image{
		Project:  project,
		RepoName: repo,
		Tag:      tag,
	}, nil
}

func runScan(instance *HarborInstance, access *HarborAccess, image *Image) {
	//TODO implement
	print(instance)
	print(access)
	print(image)
}

func writeComment(github *Github) {
	//TODO implement
	print(github)
}

func isStringEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func isStringPresent(s string) bool {
	return !isStringEmpty(s)
}

func isPortValid(portString string) bool {
	port, err := strconv.Atoi(portString)
	if err != nil {
		return false
	} else {
		const MinPort = 1
		const MaxPort = 65535
		return port >= MinPort && port <= MaxPort
	}
}
func exitOnError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
