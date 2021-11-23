package errors

import (
	"github.com/pkg/errors"
	"net/http"
)

// 创建普通错误
func New(message string) error {
	return errors.New(message)
}

// Errorf 创建新错误
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// 业务错误
type codeError struct {
	code int
	err  string
}

func (c codeError) Error() string {
	return c.err
}

// CodeError 新建业务错误，附带错误码
func CodeError(code int, err string) error {
	return codeError{code: code, err: err}
}

func BadError(err string) error {
	return codeError{code: http.StatusBadRequest, err: err}
}

// Code 提取错误码，codeError 返回 code 和 true，其他返回 0 和 false
func Code(err error) (int, bool) {
	err = errors.Cause(err)
	if ce, ok := err.(codeError); ok {
		return ce.code, true
	}
	return 0, false
}
