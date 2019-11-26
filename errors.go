package errors

import (
	"errors"
	"fmt"
)

type ErrorExt interface {
	error
	AddTag(key string, value interface{}) ErrorExt
	AddTags(tags ...Tag) ErrorExt
	AddTypes(types ...string) ErrorExt
}

// Returns a new error with no cause
func New(message string) ErrorExt {
	return &wrapper{
		nil,
		message,
		callers(),
		nil,
		nil,
	}
}

// Returns a new error with no cause
func Newf(format string, args ...interface{}) ErrorExt {
	return &wrapper{
		nil,
		fmt.Sprintf(format, args...),
		callers(),
		nil,
		nil,
	}
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) ErrorExt {
	if err == nil {
		return nil
	}
	return &wrapper{
		err,
		message,
		callers(),
		nil,
		nil,
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) ErrorExt {
	if err == nil {
		return nil
	}
	return &wrapper{
		err,
		fmt.Sprintf(format, args...),
		callers(),
		nil,
		nil,
	}
}

// Unwrap returns the wrapped error in the given error, if any. Delegates to the standard "errors" lib.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// RootCause returns the root, underlying, cause of an error, if possible.
func RootCause(err error) error {
	for err != nil {
		cause := Unwrap(err)
		if cause == nil {
			break
		} else {
			err = cause
		}
	}
	return err
}
