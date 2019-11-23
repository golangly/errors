package errors

import (
	"fmt"
)

type ErrorExt interface {
	error
	AddTag(key string, value interface{}) ErrorExt
	AddTags(tags ...Tag) ErrorExt
	AddTypes(types ...string) ErrorExt
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) ErrorExt {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
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
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
		nil,
		nil,
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) ErrorExt {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
		nil,
		nil,
	}
}
