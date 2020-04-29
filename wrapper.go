package errors

import (
	"fmt"
	"io"
	"strings"

	"github.com/kr/text"
)

type wrapper struct {
	cause   error
	message string
	stack   *stack
	types   []string
	tags    []Tag
}

// Return the causing (wrapped) error of this error.
func (w *wrapper) Error() string { return w.message }

// Return the causing (wrapped) error of this error.
func (w *wrapper) Cause() error { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *wrapper) Unwrap() error { return w.cause }

// Formats this error for printf strings.
func (w *wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		_, _ = io.WriteString(s, w.Error())
		if s.Flag('+') {
			for _, tag := range w.tags {
				_, _ = fmt.Fprintf(s, " %s=%v", tag.Key, tag.Value)
			}
			if len(w.types) > 0 {
				_, _ = fmt.Fprintf(s, " types=%s", strings.Join(w.types, ","))
			}

			indent := ""
			if precision, ok := s.Precision(); ok {
				indent = strings.Repeat(" ", precision)
			}
			sw := stateWrapper{
				w:      text.NewIndentWriter(s, []byte(indent)),
				target: s,
			}
			w.stack.Format(sw, verb)
			_, _ = fmt.Fprintf(s, "\n%+v", w.Cause())
		}
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

// Adds the given tag.
func (w *wrapper) AddTag(key string, value interface{}) ErrorExt {
	return w.AddTags(T(key, value))
}

// Adds the given tags.
func (w *wrapper) AddTags(tags ...Tag) ErrorExt {
	w.tags = append(w.tags, tags...)
	return w
}

// Adds the given types.
func (w *wrapper) AddTypes(types ...string) ErrorExt {
	w.types = append(w.types, types...)
	return w
}

type stateWrapper struct {
	w      io.Writer
	target fmt.State
}

func (i stateWrapper) Write(b []byte) (n int, err error) {
	return i.w.Write(b)
}

func (i stateWrapper) Width() (wid int, ok bool) {
	return i.target.Width()
}

func (i stateWrapper) Precision() (prec int, ok bool) {
	return i.target.Precision()
}

func (i stateWrapper) Flag(c int) bool {
	return i.target.Flag(c)
}
