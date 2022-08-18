package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/operation"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/webutil"
	"io"
	"strconv"
	"time"
)

const (
	ErrFailed        = "General Failure"
	ErrImageNotFound = "ImageNotFound"
)

type Finder struct {
}

func GetFinder() *Finder {
	return &Finder{}
}

func (f *Finder) FindImage() *operation.Result {
	image := config.Get().ImageInfo.Raw
	instance := config.Get().Harbor.Instance.Host

	var result *operation.Result
	for {
		if config.IsTimeoutSet() {
			if config.GetTimer().IsTimeOver() {
				err := fmt.Sprintf("Failed to find image after %d seconds. Got timeout! \n",
					config.Get().Timing.Timeout)
				util.ExitOnError(errors.New(err))
				return operation.Failed()
			}
		}
		result = doCheck()
		if result.Ok() {
			log.Info.Printf("Image '%s' found in '%s' instance \n", image, instance)
			break
		} else {
			if result.HasMessage() && result.GetMessage() == ErrImageNotFound {
				//no image
				log.Info.Printf("Image '%s' is not found yet. Pausing for %d seconds \n",
					image, config.Get().Timing.CheckInterval)
				doPause()
			} else if result.HasError() {
				util.ExitOnError(result.GetError())
			} else {
				//unexpected
				err := errors.New("failed to search for image. Something went wrong. Got unknown error")
				util.ExitOnError(err)
			}
		}
	}

	return result
}

func doCheck() *operation.Result {
	resp, err := webutil.DoFindRequest()
	if err != nil {
		return operation.FailedWithMessageAndError("No response from request", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte
	var goodResponse harbor.FindImageJson

	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &goodResponse)
		if err != nil {
			return operation.FailedWithError(err)
		}
		if goodResponse == nil || len(goodResponse) == 0 {
			return operation.FailedWithMessage(ErrImageNotFound)
		}
		digest := config.Get().ImageInfo.Digest
		if util.IsStringPresent(digest) {
			if digest == goodResponse[0].Digest {
				return operation.Successful()
			} else {
				log.Warning.Printf("Got digest mismatch. Expected '%s', got '%s'.", digest, goodResponse[0].Digest)
				return operation.FailedWithMessage(ErrImageNotFound)
			}
		}
		return operation.Successful()
	} else {
		switch resp.StatusCode {
		case 400:
			log.Error.Println("Got Bad Request")
			return operation.FailedWithMessage(ErrFailed)
		case 401:
			log.Error.Println("Cannot request Harbor with given credentials: unauthorized. " +
				"Please check configuration again.")
			return operation.FailedWithMessage(ErrFailed)
		case 404:
			log.Error.Println("Image not found.")
			return operation.FailedWithMessage(ErrImageNotFound)
		case 500:
			log.Error.Println("Harbor-Side error")
			return operation.FailedWithMessage(ErrFailed)
		default:
			log.Error.Println("Got unexpected status " + strconv.Itoa(resp.StatusCode))
			return operation.FailedWithMessage(ErrFailed)
		}
	}
}

func doPause() {
	interval := config.Get().Timing.CheckInterval
	time.Sleep(time.Duration(interval) * time.Second)
	config.GetTimer().SecondsLeft = config.GetTimer().SecondsLeft - interval
}
