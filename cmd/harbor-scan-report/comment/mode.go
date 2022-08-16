package comment

import (
	"errors"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"strings"
)

type Mode string

const (
	CreateNew = "create_new"
	Update    = "update"
)

func CreateCommentMode(s string) (error, Mode) {
	if util.IsStringEmpty(s) {
		return nil, CreateNew
	}
	str := strings.ToLower(s)
	switch str {
	case "update":
		return nil, Update
	case "create_new":
		return nil, CreateNew
	default:
		return errors.New("wrong comment mode"), CreateNew
	}
}
