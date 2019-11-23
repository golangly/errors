package errors

import (
	"fmt"
	"io"
	"strings"
)

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct {
	msg string
	*stack
	types []string
	tags  []Tag
}

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) ErrorExt {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) ErrorExt {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// Returns the error message of the error object.
func (f *fundamental) Error() string { return f.msg }

// Adds the given tag.
func (f *fundamental) AddTag(key string, value interface{}) ErrorExt {
	return f.AddTags(T(key, value))
}

// Adds the given tags.
func (f *fundamental) AddTags(tags ...Tag) ErrorExt {
	f.tags = append(f.tags, tags...)
	return f
}

// Adds the given types.
func (f *fundamental) AddTypes(types ...string) ErrorExt {
	f.types = append(f.types, types...)
	return f
}

// Formats the error for a printf string.
func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, f.msg)
			for _, tag := range f.tags {
				_, _ = fmt.Fprintf(s, " %s=%v", tag.Key, tag.Value)
			}
			if len(f.types) > 0 {
				_, _ = fmt.Fprintf(s, " types=%s", strings.Join(f.types, ","))
			}
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, f.msg)
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", f.msg)
	}
}
