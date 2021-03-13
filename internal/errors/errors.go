package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

type errorBase struct {
	cause   error
	message string
	frame   xerrors.Frame
}

func newErrorBase() *errorBase {
	err := &errorBase{}

	err.frame = xerrors.Caller(2)

	return err
}

func (r errorBase) Error() string {
	if r.cause != nil {
		return r.cause.Error()
	}
	return r.message
}

func (r *errorBase) Unwrap() error {
	return r.cause
}

func (r *errorBase) Format(s fmt.State, v rune) {
	xerrors.FormatError(r, s, v)
}

func (r *errorBase) FormatError(p xerrors.Printer) error {
	p.Print(r.message)
	r.frame.Format(p)
	return r.cause
}

type DataBaseAccessErrorCode uint

const (
	RecordNotFound           DataBaseAccessErrorCode = iota
	OtherDataBaseAccessError DataBaseAccessErrorCode = 99
)

type DataBaseAccessError struct {
	*errorBase
	code DataBaseAccessErrorCode
}

func NewDataBaseAccessError(code DataBaseAccessErrorCode, message string, cause error) *DataBaseAccessError {
	err := &DataBaseAccessError{errorBase: newErrorBase()}

	err.message = message
	err.cause = cause
	err.code = code

	return err
}

func (r *DataBaseAccessError) Code() DataBaseAccessErrorCode {
	return r.code
}

type TimeoutError struct {
	errorBase
}

func NewTimeoutError(msg string) TimeoutError {
	return TimeoutError{errorBase{message: msg}}
}

type ApplicationError struct {
	*errorBase
	code          string
	messageParams []interface{}
}

type ApplicationErrorOption func(*ApplicationError)

func Message(message string) ApplicationErrorOption {
	return func(err *ApplicationError) {
		err.message = message
	}
}

func Cause(cause error) ApplicationErrorOption {
	return func(err *ApplicationError) {
		err.cause = cause
	}
}

func Code(code string) ApplicationErrorOption {
	return func(err *ApplicationError) {
		err.code = code
	}
}

func MessageParams(messageParams ...interface{}) ApplicationErrorOption {
	return func(err *ApplicationError) {
		err.messageParams = messageParams
	}
}

func NewApplicationError(options ...ApplicationErrorOption) *ApplicationError {
	err := &ApplicationError{errorBase: newErrorBase()}

	for _, opt := range options {
		opt(err)
	}

	return err
}

func (r *ApplicationError) Code() string {
	return r.code
}

func (r *ApplicationError) MessageParams() []interface{} {
	return r.messageParams
}
