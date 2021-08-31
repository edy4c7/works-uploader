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

type RecordNotFoundError struct {
	*errorBase
}

func NewRecordNotFoundError(message string, cause error) *RecordNotFoundError {
	err := &RecordNotFoundError{errorBase: newErrorBase()}

	err.message = message
	err.cause = cause

	return err
}

type BadRequestError struct {
	*errorBase
}

func NewBadRequestError(message string, cause error) *BadRequestError {
	err := &BadRequestError{errorBase: newErrorBase()}

	err.message = message
	err.cause = cause

	return err
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

	err.message = "an ApplicationError has occurred"

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
