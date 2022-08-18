package operation

import "github.com/kyberorg/harbor-scan-report/cmd/harbor-scan-report/util"

type Result struct {
	successFlag bool
	err         error
	message     string
}

func (r Result) Ok() bool {
	return r.successFlag
}

func (r Result) Failed() bool {
	return !r.successFlag
}

func (r Result) HasError() bool {
	return r.err != nil
}

func (r Result) GetError() error {
	return r.err
}

func (r Result) HasMessage() bool {
	return util.IsStringPresent(r.message)
}

func (r Result) GetMessage() string {
	return r.message
}

func (r Result) SetMessage(message string) {
	r.message = message
}

func Successful() *Result {
	return &Result{
		successFlag: true,
		err:         nil,
		message:     "",
	}
}

func Failed() *Result {
	return &Result{
		successFlag: false,
		err:         nil,
		message:     "",
	}
}

func FailedWithMessage(msg string) *Result {
	return &Result{
		successFlag: false,
		err:         nil,
		message:     msg,
	}
}

func FailedWithError(err error) *Result {
	return &Result{
		successFlag: false,
		err:         err,
		message:     "",
	}
}

func FailedWithMessageAndError(msg string, err error) *Result {
	return &Result{
		successFlag: false,
		err:         err,
		message:     msg,
	}
}
