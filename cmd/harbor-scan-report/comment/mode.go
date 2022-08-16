package comment

import (
	"errors"
	"fmt"
	"github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"
	"strings"
)

type Mode string

const (
	CreateNew  = "create_new"
	UpdateLast = "update_last"
)

func CreateCommentMode(s string) (error, Mode) {
	if util.IsStringEmpty(s) {
		return nil, CreateNew
	}
	str := strings.ToLower(s)
	switch str {
	case "update_last", "update":
		return nil, UpdateLast
	case "new", "create_new":
		return nil, CreateNew
	default:
		return errors.New(fmt.Sprintf("'%s' is wrong comment mode. Valid values: 'create_new' and 'update_last'", str)), CreateNew
	}
}
