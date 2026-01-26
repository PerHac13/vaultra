package app

import (
	"fmt"
)

type ErrorKind string

const (
	ErrKindDatabase	    ErrorKind = "DatabaseError"
	ErrKindStorage      ErrorKind = "StorageError"
	ErrKindCompression  ErrorKind = "CompressionError"
	ErrKindConfig       ErrorKind = "ConfigError"
	ErrKindValidation   ErrorKind = "ValidationError"
)

type Error struct {
	Op  string
	Kind ErrorKind
	Err error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Kind, e.Err)
	}

	return fmt.Sprintf("%s: %s", e.Op, e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewError(op string, kind ErrorKind, err error) *Error {
	return &Error{
		Op:   op,
		Kind: kind,
		Err:  err,
	}
}