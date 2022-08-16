package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/comment"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/config"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/harbor"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/log"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/scan"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/severity"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/webutil"
	"io"
	"net/http"
	"strings"
)

const ErrorMarker = "ErrorMarker"
const CommentNotFoundMarker = "CommentNotFound"

var report *scan.Report

func WriteComment(scanReport *scan.Report) {
	report = scanReport
	message := createMessage()
	commentMode := config.Get().Comment.Mode
	var resp *http.Response
	var err error
	if commentMode == comment.Update {
		//search for existing comment
		var commentUpdateUrl string
		commentUpdateUrl, err = searchForExistingComment()
		switch commentUpdateUrl {
		case ErrorMarker:
			log.Warning.Println("Failed to update previous comment: got error while searching")
			log.Debug.Printf("search error: " + err.Error())
			resp, err = webutil.DoGitHubCommentCreateRequest(message)
			break
		case CommentNotFoundMarker:
			log.Warning.Println("Failed to update previous comment: comment not found")
			resp, err = webutil.DoGitHubCommentCreateRequest(message)
			break
		default:
			//all good - updating comment
			resp, err = webutil.DoGitHubCommentUpdateRequest(commentUpdateUrl, message)
		}
	} else {
		resp, err = webutil.DoGitHubCommentCreateRequest(message)
	}

	if err != nil {
		log.Warning.Printf("Failed to create GitHub Comment")
		util.ExitOnError(err)
	}
	if resp.StatusCode == 200 {
		log.Info.Println("GitHub comment updated")
	} else if resp.StatusCode == 201 {
		log.Info.Println("GitHub comment created")
	} else {
		log.Warning.Printf("Failed to create GitHub comment. Status: %d \n", resp.StatusCode)
	}
}

func createMessage() string {
	var b strings.Builder

	b.WriteString(getTitle())
	b.WriteString(fmt.Sprintf("Results for image [%s](%s) \n", config.Get().ImageInfo.Raw, harbor.UiUrl()))
	b.WriteString(topSeverityEmoji() + " ")
	b.WriteString(fmt.Sprintf("Total %d vulnerabilities found ",
		report.Counters.Total))
	if report.Counters.Total > 0 {
		b.WriteString(fmt.Sprintf("- %d fixable ", report.Counters.Fixable))
	}
	b.WriteString(fmt.Sprintf("\n"))
	if report.Counters.Total > 0 {
		b.WriteString(fmt.Sprintf(
			"[%s](## \"critical\") %d critical "+
				"[%s](## \"high\") %d high "+
				"[%s](## \"medium\") %d medium "+
				"[%s](## \"low\") %d low\n",
			s2e(severity.Critical), report.Counters.Critical,
			s2e(severity.High), report.Counters.High,
			s2e(severity.Medium), report.Counters.Medium,
			s2e(severity.Low), report.Counters.Low,
		))
	}
	b.WriteString(fmt.Sprintf("Scanned with `%s %s` from `%s` \n",
		report.Scanner.Name, report.Scanner.Version, report.Scanner.Vendor))
	b.WriteString(fmt.Sprintf("Report generated at `%s`\n", util.PrettyDate(report.GeneratedAt)))

	return b.String()
}

func getTitle() string {
	return fmt.Sprintf("## %s \n", config.Get().Comment.Title)
}

func s2e(s severity.Severity) string {
	switch s {
	case severity.Critical:
		return ":no_entry:"
	case severity.High:
		return ":fire:"
	case severity.Medium:
		return ":warning:"
	case severity.Low:
		return ":triangular_flag_on_post:"
	case severity.None:
		return ":heavy_check_mark:"
	default:
		return ":interrobang:"
	}
}

func topSeverityEmoji() string {
	return s2e(report.TopSeverity)
}

func searchForExistingComment() (string, error) {
	resp, err := webutil.DoGitHubCommentSearchRequest()
	if err != nil {
		return ErrorMarker, err
	}
	if resp.StatusCode == 200 {
		var issueComments IssueComments
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrorMarker, err
		}
		err = json.Unmarshal(body, &issueComments)
		if err != nil {
			return ErrorMarker, err
		}
		commentUpdateUrl := CommentNotFoundMarker
		for _, c := range issueComments {
			if strings.HasPrefix(c.Body, getTitle()) {
				commentUpdateUrl = c.URL
			}
		}
		return commentUpdateUrl, nil
	} else {
		switch resp.StatusCode {
		case 404:
			return ErrorMarker, errors.New("search for comments failed - no such issue found")
		case 410:
			return ErrorMarker, errors.New("search for comments failed - issue is gone")
		default:
			return ErrorMarker, errors.New("search for comments failed - unknown response code")
		}
	}
}
