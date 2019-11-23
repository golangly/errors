package errors

import (
	"fmt"
	"io"
	"strings"
)

type withMessage struct {
	cause error
	msg   string
	types []string
	tags  []Tag
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }

func (w *withMessage) Cause() error { return w.cause }

func (w *withMessage) AddTag(key string, value interface{}) ErrorExt { return w.AddTags(T(key, value)) }
func (w *withMessage) AddTags(tags ...Tag) ErrorExt {
	w.tags = append(w.tags, tags...)
	return w
}
func (w *withMessage) AddTypes(types ...string) ErrorExt {
	w.types = append(w.types, types...)
	return w
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", w.Cause())
			_, _ = io.WriteString(s, w.msg)
			for _, tag := range w.tags {
				_, _ = fmt.Fprintf(s, " %s=%v", tag.Key, tag.Value)
			}
			if len(w.types) > 0 {
				_, _ = fmt.Fprintf(s, " types=%s", strings.Join(w.types, ","))
			}
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, w.Error())
	}
}
